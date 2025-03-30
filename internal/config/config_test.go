package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/brpaz/freshrss-cleaner/internal/config"
)

func TestGetDefaultConfigFilePath(t *testing.T) {
	filePath := config.DefaultConfigFilePath()
	assert.NotEmpty(t, filePath)

	configDir, err := os.UserConfigDir()
	require.NoError(t, err)
	expected := filepath.Join(configDir, "freshrss-cleaner.yaml")
	assert.Equal(t, expected, filePath)
}

func TestCreateDefaultConfigFile(t *testing.T) {
	t.Run("Returns Path if config file already exists", func(t *testing.T) {
		configFilePath := t.TempDir()
		configFilePath = filepath.Join(configFilePath, "freshrss-cleaner.yaml")

		// Create the file
		err := os.WriteFile(configFilePath, []byte("test"), 0o600)
		require.NoError(t, err)

		defer os.Remove(configFilePath)

		path, err := config.CreateDefaultConfigFile(configFilePath)
		require.NoError(t, err)

		assert.Equal(t, configFilePath, path)
	})

	t.Run("Creates a config file in the specified location", func(t *testing.T) {
		configFilePath := t.TempDir()
		configFilePath = filepath.Join(configFilePath, "freshrss-cleaner.yaml")

		createdFilePath, err := config.CreateDefaultConfigFile(configFilePath)
		require.NoError(t, err)

		// Clean up
		defer os.Remove(createdFilePath)

		assert.Equal(t, configFilePath, createdFilePath)

		// Check if the file was created
		_, err = os.Stat(createdFilePath)
		assert.NoError(t, err)
	})
}
