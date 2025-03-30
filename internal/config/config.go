// Package config provides functionality to manage the configuration of the application.
package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// RootConfig represents the root configuration structure for the application.
type RootConfig struct {
	URL      string       `yaml:"url"`
	Username string       `yaml:"username"`
	Password string       `yaml:"password"`
	Feeds    []FeedConfig `yaml:"feeds"`
}

// FeedConfig represents the configuration for a specific feed.
type FeedConfig struct {
	ID   string `yaml:"id"`
	Days int    `yaml:"days"`
}

// DefaultConfig provides the default configuration template
const DefaultConfig = `url: "https://<your-freshrss-instance>"
username: "user"
password: "pass"
api_key: ""
feeds:
  - id: "feed1"
	days: 7
`

// DefaultConfigFilePath returns the default path for the configuration file.
func DefaultConfigFilePath() string {
	configDir, _ := os.UserConfigDir()
	return filepath.Join(configDir, filename)
}

// CreateDefaultConfigFile creates a default configuration file in the user's config directory.
func CreateDefaultConfigFile(configFilePath string) (string, error) {
	// Check if config file already exists
	if _, err := os.Stat(configFilePath); err == nil {
		return configFilePath, nil
	}

	// Ensure the directory exists
	configDir := filepath.Dir(configFilePath)
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write the config file
	if err := os.WriteFile(configFilePath, []byte(DefaultConfig), 0o600); err != nil {
		return "", fmt.Errorf("failed to create default config file: %w", err)
	}

	return configFilePath, nil
}
