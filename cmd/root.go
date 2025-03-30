// Package cmd contains the command definitions for the application.
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/brpaz/freshrss-cleaner/cmd/clean"
	"github.com/brpaz/freshrss-cleaner/cmd/createconfig"
	"github.com/brpaz/freshrss-cleaner/cmd/version"
)

// NewRootCmd returns a new instance of the root command for the application
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "freshrss-cleaner",
		Short: "A command line tool to clean up old entries from FreshRSS",
	}

	// Reggister subcommands
	rootCmd.AddCommand(version.New())
	rootCmd.AddCommand(clean.New())
	rootCmd.AddCommand(createconfig.New())

	return rootCmd
}
