package freshrss

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/brpaz/freshrss-cleaner/internal/config"
)

// API defines the interface for the FreshRSS API client
type API interface {
	GetAuthToken(ctx context.Context) (string, error)
	MarkAsRead(ctx context.Context, authToken string, feedID string, days int) error
}

// Cleaner is a struct that represents a Freshrss cleaner
type Cleaner struct {
	client API
	config *config.RootConfig
}

// Validate checks if the Freshrss cleaner is properly configured with all the required fields
func (c *Cleaner) Validate() error {
	if c.client == nil {
		return fmt.Errorf("client is required")
	}

	if c.config == nil {
		return fmt.Errorf("config is required")
	}

	return nil
}

// WithClient provides a way to set the Freshrss client for the cleaner
func WithClient(client API) CleanerOption {
	return func(c *Cleaner) {
		c.client = client
	}
}

func WithConfig(cfg *config.RootConfig) CleanerOption {
	return func(c *Cleaner) {
		c.config = cfg
	}
}

// Option defines a function to configure the FreshRSS client
type CleanerOption func(*Cleaner)

// NewCleaner creates a new instance of the Freshrss cleaner
func NewCleaner(opts ...CleanerOption) (*Cleaner, error) {
	cleaner := &Cleaner{}

	for _, opt := range opts {
		opt(cleaner)
	}

	if err := cleaner.Validate(); err != nil {
		return nil, fmt.Errorf("invalid cleaner configuration: %w", err)
	}

	return cleaner, nil
}

// CleanOldEntries cleans up old entries from FreshRSS based on the provided configuration
func (c *Cleaner) CleanOldEntries(ctx context.Context, log *slog.Logger) error {
	log.Info("Fetching auth token")
	authToken, err := c.client.GetAuthToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get auth token: %w", err)
	}

	for _, feed := range c.config.Feeds {
		log.Info("Processing feed", "feed_id", feed.ID)
		err := c.processFeed(ctx, feed, authToken)
		if err != nil {
			log.Error("Failed to process feed", "feed_id", feed.ID, "error", err)
			continue
		}

	}

	return nil
}

func (c *Cleaner) processFeed(ctx context.Context, feed config.FeedConfig, authToken string) error {
	return c.client.MarkAsRead(ctx, authToken, feed.ID, feed.Days)
}
