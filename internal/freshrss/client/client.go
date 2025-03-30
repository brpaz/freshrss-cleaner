// Package client provides a Go client for the FreshRSS API
package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client represents a client for the FreshRSS API
type Client struct {
	baseURL    string
	username   string
	password   string
	httpClient *http.Client
}

// Validate checks if the client is configured properly
func (c *Client) Validate() error {
	if c.baseURL == "" {
		return fmt.Errorf("base URL is required")
	}

	if _, err := url.ParseRequestURI(c.baseURL); err != nil {
		return fmt.Errorf("invalid base URL: %w", err)
	}

	if c.username == "" {
		return fmt.Errorf("username is required")
	}

	if c.password == "" {
		return fmt.Errorf("password is required")
	}

	return nil
}

// Option defines a function to configure the FreshRSS client
type Option func(*Client)

// WithBaseURL sets the base URL for the FreshRSS API
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		// Remove trailing slashes for consistency
		c.baseURL = strings.TrimRight(baseURL, "/")
	}
}

// WithCredentials sets the username and password for authentication
func WithCredentials(username, password string) Option {
	return func(c *Client) {
		c.username = username
		c.password = password
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithTimeout sets a custom timeout for the HTTP client
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		if c.httpClient == nil {
			c.httpClient = &http.Client{}
		}
		c.httpClient.Timeout = timeout
	}
}

// New creates a new FreshRSS client with the provided options
func New(opts ...Option) (*Client, error) {
	client := &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	for _, opt := range opts {
		opt(client)
	}

	if err := client.Validate(); err != nil {
		return nil, fmt.Errorf("invalid FreshRSS client configuration: %w", err)
	}

	return client, nil
}

// setAuthHeaders adds authentication headers to an HTTP request
func (c *Client) setAuthHeaders(req *http.Request, authToken string) {
	if authToken != "" {
		req.Header.Set("Authorization", "GoogleLogin auth="+authToken)
	}
}

// GetAuthToken retrieves an authentication token from the FreshRSS API
func (c *Client) GetAuthToken(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/accounts/ClientLogin", nil)
	if err != nil {
		return "", fmt.Errorf("error creating auth request: %w", err)
	}

	query := req.URL.Query()
	query.Add("Email", c.username)
	query.Add("Passwd", c.password)
	req.URL.RawQuery = query.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error executing auth request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading auth response body: %w", err)
	}

	// Parse the response to get the auth token
	lines := strings.Split(string(body), "\n")
	if len(lines) < 3 {
		return "", fmt.Errorf("unexpected auth response format")
	}

	for _, line := range lines {
		if strings.HasPrefix(line, "Auth=") {
			return strings.TrimPrefix(line, "Auth="), nil
		}
	}

	return "", fmt.Errorf("auth token not found in response")
}

// MarkAsRead marks items in a feed as read that are older than the specified days
func (c *Client) MarkAsRead(ctx context.Context, authToken string, feedID string, olderThanDays int) error {
	if authToken == "" {
		return fmt.Errorf("auth token is required")
	}

	if feedID == "" {
		return fmt.Errorf("feed ID is required")
	}

	// Calculate cutoff time
	cutoffTime := time.Now().AddDate(0, 0, -olderThanDays).UnixNano() / 1e3 // Microseconds

	// Prepare request
	endpoint := fmt.Sprintf("%s/reader/api/0/mark-all-as-read", c.baseURL)

	// Prepare form data
	data := url.Values{}
	data.Set("s", feedID)
	data.Set("ts", fmt.Sprintf("%d", cutoffTime))

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return fmt.Errorf("error creating mark-as-read request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	c.setAuthHeaders(req, authToken)

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error executing mark-as-read request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("mark-as-read request failed with unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
