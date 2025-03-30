// Package clean provides the command definition for the clean command.
package clean

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/brpaz/freshrss-cleaner/internal/config"
	"github.com/brpaz/freshrss-cleaner/internal/freshrss"
	"github.com/brpaz/freshrss-cleaner/internal/freshrss/client"
)

// New creates a new clean command
func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Clean up old entries from FreshRSS",
		RunE:  runClean,
	}

	cmd.Flags().StringP("config", "c", config.DefaultConfigFilePath(), "Path to the configuration file")

	return cmd
}

// runClean handles the execution of the clean command
func runClean(cmd *cobra.Command, args []string) error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Get configuration path
	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return fmt.Errorf("failed to get config flag: %w", err)
	}

	// Load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config file: %w", err)
	}

	// Initialize FreshRSS client
	client, err := createFreshRSSClient(cfg)
	if err != nil {
		return err
	}

	// Run the cleaner
	cleaner, err := freshrss.NewCleaner(
		freshrss.WithClient(client),
		freshrss.WithConfig(cfg),
	)
	if err != nil {
		return fmt.Errorf("failed to create cleaner: %w", err)
	}

	ctx := cmd.Context()
	if err := cleaner.CleanOldEntries(ctx, logger); err != nil {
		return fmt.Errorf("failed to run cleaner: %w", err)
	}

	return nil
}

// createFreshRSSClient initializes a new FreshRSS client with configuration
func createFreshRSSClient(cfg *config.RootConfig) (*client.Client, error) {
	client, err := client.New(
		client.WithBaseURL(cfg.URL),
		client.WithCredentials(cfg.Username, cfg.Password),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create freshrss client: %w", err)
	}

	return client, nil
}
