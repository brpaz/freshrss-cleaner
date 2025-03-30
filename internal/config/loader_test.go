package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/brpaz/freshrss-cleaner/internal/config"
)

func TestLoad(t *testing.T) {
	t.Run("With non existing config file", func(t *testing.T) {
		_, err := config.Load("nonexistent.yaml")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no such file or directory")
	})

	t.Run("With empty config file", func(t *testing.T) {
		_, err := config.Load("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config path cannot be empty")
	})

	t.Run("With valid config file", func(t *testing.T) {
		configFile := "testdata/valid_config.yaml"
		cfg, err := config.Load(configFile)
		require.NoError(t, err)
		assert.Equal(t, "https://example.com", cfg.URL)
		assert.Equal(t, "user", cfg.Username)
		assert.Equal(t, "pass", cfg.Password)
		assert.Len(t, cfg.Feeds, 2)
		assert.Equal(t, "feed1", cfg.Feeds[0].ID)
		assert.Equal(t, 7, cfg.Feeds[0].Days)
		assert.Equal(t, "feed2", cfg.Feeds[1].ID)
		assert.Equal(t, 14, cfg.Feeds[1].Days)
	})

	t.Run("With invalid config file", func(t *testing.T) {
		configFile := "testdata/invalid_config.yaml"
		_, err := config.Load(configFile)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse config file")
	})

	t.Run("With valid config file with env var replacement", func(t *testing.T) {
		configFile := "testdata/valid_config_with_env.yaml"

		t.Setenv("FRESHRSS_URL", "https://example.com")
		t.Setenv("FRESHRSS_USERNAME", "user")
		t.Setenv("FRESHRSS_PASSWORD", "pass")

		cfg, err := config.Load(configFile)
		require.NoError(t, err)

		assert.Equal(t, "https://example.com", cfg.URL)
		assert.Equal(t, "user", cfg.Username)
		assert.Equal(t, "pass", cfg.Password)
	})
}
