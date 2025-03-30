package freshrss_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/brpaz/freshrss-cleaner/internal/config"
	"github.com/brpaz/freshrss-cleaner/internal/freshrss"
)

// mockClient implements a mock of the client interface for testing
type mockClient struct {
	mock.Mock
}

func (m *mockClient) GetAuthToken(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *mockClient) MarkAsRead(ctx context.Context, authToken string, feedID string, days int) error {
	args := m.Called(ctx, authToken, feedID, days)
	return args.Error(0)
}

// Test fixtures
var mockConfig = &config.RootConfig{
	URL:      "https://example.com",
	Username: "user",
	Password: "pass",
	Feeds: []config.FeedConfig{
		{ID: "feed1", Days: 7},
		{ID: "feed2", Days: 14},
	},
}

func TestNewCleaner(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		client      freshrss.API
		config      *config.RootConfig
		expectError bool
		errorMsg    string
	}{
		{
			name:        "WithValidConfigs",
			client:      &mockClient{},
			config:      mockConfig,
			expectError: false,
		},
		{
			name:        "WithInvalidClient",
			client:      nil,
			config:      mockConfig,
			expectError: true,
			errorMsg:    "client is required",
		},
		{
			name:        "WithInvalidConfig",
			client:      &mockClient{},
			config:      nil,
			expectError: true,
			errorMsg:    "config is required",
		},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c, err := freshrss.NewCleaner(
				freshrss.WithClient(tc.client),
				freshrss.WithConfig(tc.config),
			)

			if tc.expectError {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
				assert.Nil(t, c)
			} else {
				assert.Nil(t, err)
				assert.IsType(t, &freshrss.Cleaner{}, c)
			}
		})
	}
}

func TestCleanOldEntries(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		client := &mockClient{}
		ctx := context.Background()

		cleaner, err := freshrss.NewCleaner(
			freshrss.WithClient(client),
			freshrss.WithConfig(mockConfig),
		)
		assert.Nil(t, err)

		client.On("GetAuthToken", ctx).Return("mockToken", nil)
		client.On("MarkAsRead", ctx, "mockToken", "feed1", 7).Return(nil)
		client.On("MarkAsRead", ctx, "mockToken", "feed2", 14).Return(nil)

		err = cleaner.CleanOldEntries(ctx, logger)
		assert.Nil(t, err)

		client.AssertExpectations(t)
	})

	t.Run("Returns error when GetAuthToken fails", func(t *testing.T) {
		t.Parallel()
		client := &mockClient{}
		ctx := context.Background()

		cleaner, err := freshrss.NewCleaner(
			freshrss.WithClient(client),
			freshrss.WithConfig(mockConfig),
		)
		assert.Nil(t, err)

		client.On("GetAuthToken", ctx).Return("", assert.AnError)

		err = cleaner.CleanOldEntries(ctx, logger)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "failed to get auth token")

		client.AssertExpectations(t)
	})
}
