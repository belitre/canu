package uninstall

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/belitre/canu/pkg/config"
	"github.com/belitre/canu/pkg/utils"
	"github.com/spf13/cobra"
)

func UninstallCmd(c *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "uninstalls canu aliases",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("error, this command doesn't accept arguments")
			}

			return start(c)
		},
	}

	defaultShellConfigPath := path.Join(c.UserHomePath, config.DefaultShellConfigScript)

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

	fmt.Printf("trying to remove file %s for alias %s...\n", scriptPath, cfg.AliasName)

	if err := os.Remove(scriptPath); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("warning: script %s for alias %s not found. Uninstall will continue...\n", scriptPath, cfg.AliasName)
		} else {
			fmt.Printf("error while removing file %s for alias %s. You may need to remove the file manually...\n", scriptPath, cfg.AliasName)
		}
	} else {
		fmt.Printf("file %s for alias %s removed successfully...\n", scriptPath, cfg.AliasName)
	}

	fmt.Printf("trying to remove alias %s from shell config script: %s ...\n", cfg.AliasName, cfg.ShellConfigScriptPath)

	aliasCommand := fmt.Sprintf("alias %s=", cfg.AliasName)

	if err := utils.RemoveAliasFromFile(cfg.ShellConfigScriptPath, aliasCommand); err != nil {
		return err
	}

	fmt.Printf("uninstall for alias %s done!\n", cfg.AliasName)

	return nil
}
