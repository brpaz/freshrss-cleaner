package config

import (
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

const filename = "freshrss-cleaner.yaml"

// envVarPattern defines the pattern for environment variable placeholders in config files
const envVarPattern = `env\("([^"]+)"\)`

var envVarRegex = regexp.MustCompile(envVarPattern)

// replaceEnv replaces environment variables in the given data with their respective values.
// Example: env("VAR_NAME") will be replaced with the value of os.Getenv("VAR_NAME").
// If the environment variable is not defined, it will be replaced with an empty string.
func replaceEnv(data []byte) []byte {
	content := string(data)

	result := envVarRegex.ReplaceAllStringFunc(content, func(match string) string {
		matches := envVarRegex.FindStringSubmatch(match)
		if len(matches) < 2 {
			return match
		}

		envVarName := matches[1]
		envVarValue := os.Getenv(envVarName)

		return envVarValue
	})

	return []byte(result)
}

// Load loads the configuration file from the specified path and returns a RootConfig struct.
// The configuration file should be in YAML format.
// Environment variables specified in the config file with env("VAR_NAME") will be replaced with their values.
// Returns an error if the config file cannot be read or parsed.
func Load(configPath string) (*RootConfig, error) {
	if configPath == "" {
		return nil, fmt.Errorf("config path cannot be empty")
	}

	// Read the configuration file
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	// Replace environment variables in the config data
	configData = replaceEnv(configData)

	// Parse the configuration file
	var config RootConfig
	if err = yaml.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
	}

	return &config, nil
}
