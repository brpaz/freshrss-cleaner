// Package main provides the entry point for the application.
// It is responsible for initializing the global dependencies and executing the root command.
package main

import (
	"fmt"
	"os"

	"github.com/brpaz/freshrss-cleaner/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
