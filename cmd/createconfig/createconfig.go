// Package createconfig provides a command to create a default configuration file for the application.
package createconfig

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/brpaz/freshrss-cleaner/internal/config"
)

// New creates a instance of the create-config command
func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-config",
		Short: "Clean a base configuration file in the user's config directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()

			configPath, err := config.CreateDefaultConfigFile(config.DefaultConfigFilePath())
			if err != nil {
				return fmt.Errorf("failed to create default config file: %w", err)
			}

			fmt.Fprintf(out, "config file created at: %s\n", configPath)

			return nil
		},
	}

	return cmd
}
