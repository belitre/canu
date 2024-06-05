package install

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/belitre/canu/pkg/config"
	"github.com/belitre/canu/pkg/utils"
	"github.com/spf13/cobra"
)

const (
	scriptContentTemplate = `
	#!/bin/sh

	%s
	
	selected_profile="$(cat %s)"
	
	if [ ! -z "$selected_profile" ]
	then
	  export AWS_PROFILE="$selected_profile"
	else
	  echo "error, %s file doesn't contain a profile name"
	fi
	`
)

func InstallCmd(c *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "installs canu aliases",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("error, this command doesn't accept arguments")
			}

			return start(c)
		},
	}

	defaultShellConfigPath := path.Join(c.UserHomePath, config.DefaultShellConfigScript)

	cmd.Flags().BoolVar(&c.IsSkipAlias, "skip-alias", false, "if provided it won't add an alias to the shell config script")
	cmd.Flags().StringVar(&c.ShellConfigScriptPath, "shell-config", defaultShellConfigPath, "path to the shell config script")
	cmd.Flags().StringVar(&c.AliasName, "alias-name", config.DefaultAliasName, "name of the alias")

	return cmd
}

func start(cfg *config.Config) error {
	cfg.AliasName = strings.TrimSpace(cfg.AliasName)

	if len(cfg.AliasName) == 0 {
		return fmt.Errorf("an alias name is required")
	}

	executablePath, err := os.Executable()

	if err != nil {
		return fmt.Errorf("error while reading path for canu binary: %v", err)
	}

	scriptPath := utils.GetScriptPath(executablePath, cfg.AliasName)

	canuConfigFile := path.Join(cfg.UserHomePath, config.CanuSaveProfileFileName)

	// TODO: build canu flags (includes, excludes, etc...)
	command := buildCommand(executablePath, cfg)

	scriptContent := fmt.Sprintf(scriptContentTemplate, command, canuConfigFile, canuConfigFile)

	fmt.Printf("creating script %s for alias %s ...\n", scriptPath, cfg.AliasName)

	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0700); err != nil {
		return fmt.Errorf("error while creating script %s: %v", scriptPath, err)
	}

	fmt.Printf("created script %s for alias %s ...\n", scriptPath, cfg.AliasName)

	aliasCommand := fmt.Sprintf("alias %s=\"source %s\"", cfg.AliasName, scriptPath)

	if cfg.IsSkipAlias {
		fmt.Printf("--skip-alias flag was provided, please add manually to your ~/.bash_profile:\n%s\n", aliasCommand)
		return nil
	}

	fmt.Printf("creating alias %s in shell config file %s ...\n", cfg.AliasName, cfg.ShellConfigScriptPath)

	aliasBlock := fmt.Sprintf("\n%s\n", aliasCommand)

	f, err := os.OpenFile(cfg.ShellConfigScriptPath, os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		return fmt.Errorf("error while opening file %s: %v", cfg.ShellConfigScriptPath, err)
	}

	defer f.Close()

	if _, err := f.WriteString(aliasBlock); err != nil {
		return fmt.Errorf("error while adding alias to file %s: %v", cfg.ShellConfigScriptPath, err)
	}

	fmt.Printf("alias %s added to %s, reload your shell or run manually: %s\n", cfg.AliasName, cfg.ShellConfigScriptPath, aliasCommand)

	return nil
}

func buildCommand(execPath string, cfg *config.Config) string {
	args := ""

	if len(cfg.ConfigPath) > 0 {
		args = fmt.Sprintf("--config-path %s", cfg.ConfigPath)
	}

	for _, include := range cfg.Include {
		args = fmt.Sprintf("%s --include %s", args, include)
	}

	for _, exclude := range cfg.Exclude {

		args = fmt.Sprintf("%s --exclude %s", args, exclude)
	}

	command := fmt.Sprintf("%s %s", execPath, args)

	fmt.Printf("command: %s\n", command)

	return command
}
