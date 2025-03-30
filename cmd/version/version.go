// Package version provides the command definition for the version command.
package version

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	BuildDate = ""
	Version   = "dev"
	GitCommit = ""
)

// New creates a command that prints the version of the application.
func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()

			fmt.Fprintf(out, "Build date: %s\n", BuildDate)
			fmt.Fprintf(out, "Version: %s\n", Version)
			fmt.Fprintf(out, "Git commit: %s\n", GitCommit)
			fmt.Fprintf(out, "Go version: %s\n", runtime.Version())

			return nil
		},
	}
	return cmd
}
