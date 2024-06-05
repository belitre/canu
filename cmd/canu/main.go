package main

import (
	"fmt"
	"os"
	"path"

	"github.com/belitre/canu/cmd/canu/install"
	"github.com/belitre/canu/cmd/canu/uninstall"
	"github.com/belitre/canu/pkg/canu"
	"github.com/belitre/canu/pkg/config"
	"github.com/belitre/canu/pkg/version"
	"github.com/spf13/cobra"
)

func main() {
	c := config.Config{}

	rootCmd := &cobra.Command{
		Use:   config.AppName,
		Short: "A CLI to switch aws profiles",
		Long:  "A CLI to switch aws profiles" + "\n" + version.GetVersion(),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("error, this command doesn't accept arguments")
			}

			return start(&c)
		},
	}

	userHome, err := os.UserHomeDir()

	if err != nil {
		fmt.Printf("error while getting value for env var $HOME, please configure your $HOME variable before runnig canu :)\n")
		os.Exit(1)
	}

	defaultPath := path.Join(userHome, config.DefaultRelativePath)
	c.UserHomePath = userHome

	installCommand := install.InstallCmd(&c)

	addCommonFlags(rootCmd, &c, defaultPath)
	addCommonFlags(installCommand, &c, defaultPath)

	rootCmd.AddCommand(installCommand)
	rootCmd.AddCommand(uninstall.UninstallCmd(&c))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func addCommonFlags(cmd *cobra.Command, c *config.Config, defaultPath string) {
	cmd.Flags().StringVarP(&c.ConfigPath, "config-path", "c", defaultPath, "Path with the AWS config file with the defined profiles")
	cmd.Flags().StringArrayVarP(&c.Include, "include", "i", []string{}, "Profiles containing any of the includes will be available to select (ignores upper/lower case)")
	cmd.Flags().StringArrayVarP(&c.Exclude, "exclude", "e", []string{}, "Profiles containing any of the excludes will be ignored. Exclude takes preference over include (ignores upper/lower case)")
}

func start(cfg *config.Config) error {
	canu := canu.New(cfg)
	return canu.Run()
}
