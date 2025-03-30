package client_test

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"

	"github.com/brpaz/freshrss-cleaner/internal/freshrss/client"
)

func initTestClient(t *testing.T) *client.Client {
	t.Helper()
	c, err := client.New(
		client.WithBaseURL("https://freshrss.example.com"),
		client.WithCredentials("test", "pass"),
		client.WithTimeout(1*time.Second),
	)

	require.NoError(t, err)
	return c
}

func mockCuttOffTime(t *testing.T, days int) int64 {
	t.Helper()
	return time.Now().AddDate(0, 0, -days).UnixNano() / 1e3 // Microseconds
}

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("WithValidConfigs", func(t *testing.T) {
		t.Parallel()
		c, err := client.New(
			client.WithBaseURL("https://example.com"),
			client.WithCredentials("user", "pass"),
			client.WithTimeout(30*time.Second),
			client.WithHTTPClient(http.DefaultClient),
		)

		assert.Nil(t, err)
		assert.IsType(t, &client.Client{}, c)
	})

	t.Run("WithInvalidBaseURL", func(t *testing.T) {
		t.Parallel()
		c, err := client.New(
			client.WithBaseURL("invalid-url"),
		)

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid base URL")
		assert.Nil(t, c)
	})

	t.Run("WithMissingUsername", func(t *testing.T) {
		t.Parallel()
		c, err := client.New(
			client.WithBaseURL("https://example.com"),
		)

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "username is required")
		assert.Nil(t, c)
	})
}

func TestGetAuthToken(t *testing.T) {
	t.Run("WithValidResponse", func(t *testing.T) {
		defer gock.Off() // Flush pending mocks after test execution

		response, err := os.ReadFile("testdata/get_auth_token_200.txt")
		require.NoError(t, err)

		// Mock the request
		gock.New("https://freshrss.example.com").
			Get("/accounts/ClientLogin").
			MatchParam("Email", "test").
			MatchParam("Passwd", "pass").
			Reply(200).
			BodyString(string(response))

		c := initTestClient(t)

		token, err := c.GetAuthToken(context.Background())
		require.NoError(t, err)
		assert.Equal(t, "test/auth-token", token)
		assert.True(t, gock.IsDone())
	})

	t.Run("WithUnauthorizedResponse", func(t *testing.T) {
		defer gock.Off() // Flush pending mocks after test execution

		response, err := os.ReadFile("testdata/get_auth_token_401.txt")
		require.NoError(t, err)

		// Mock the request
		gock.New("https://freshrss.example.com").
			Get("/accounts/ClientLogin").
			MatchParam("Email", "test").
			MatchParam("Passwd", "pass").
			Reply(401).
			BodyString(string(response))

		c := initTestClient(t)

		token, err := c.GetAuthToken(context.Background())
		require.Error(t, err)
		assert.Empty(t, token)
		assert.Equal(t, err.Error(), "request failed with status code 401: Unauthorized!\n")
		assert.True(t, gock.IsDone())
	})

	t.Run("WithInvalidResponseBody", func(t *testing.T) {
		defer gock.Off() // Flush pending mocks after test execution

		response, err := os.ReadFile("testdata/get_auth_token_invalid_response.txt")
		require.NoError(t, err)

		// Mock the request
		gock.New("https://freshrss.example.com").
			Get("/accounts/ClientLogin").
			MatchParam("Email", "test").
			MatchParam("Passwd", "pass").
			Reply(200).
			BodyString(string(response))

		c := initTestClient(t)

		token, err := c.GetAuthToken(context.Background())
		require.Error(t, err)
		assert.Empty(t, token)
		assert.Equal(t, err.Error(), "unexpected auth response format")
		assert.True(t, gock.IsDone())
	})
}

func TestMarkAsRead(t *testing.T) {
	t.Run("WithEmptyFeedID_ReturnsError", func(t *testing.T) {
		c := initTestClient(t)

		err := c.MarkAsRead(context.Background(), "test/auth-token", "", 7)
		require.Error(t, err)
		assert.Equal(t, err.Error(), "feed ID is required")
	})

	t.Run("WithEmptyAuthToken_ReturnsError", func(t *testing.T) {
		c := initTestClient(t)
		err := c.MarkAsRead(context.Background(), "", "feed-id", 7)
		require.Error(t, err)
		assert.Equal(t, err.Error(), "auth token is required")
	})

	t.Run("WithUnexpectedResponse", func(t *testing.T) {
		defer gock.Off() // Flush pending mocks after test execution

		// Mock the request
		gock.New("https://freshrss.example.com").
			Post("/reader/api/0/mark-all-as-read").
			MatchHeader("Authorization", "GoogleLogin auth=test/auth-token").
			MatchHeader("Content-Type", "application/x-www-form-urlencoded").
			Reply(401)

		c := initTestClient(t)

		err := c.MarkAsRead(context.Background(), "test/auth-token", "feed-id", 7)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "mark-as-read request failed with unexpected status code 401")
		assert.True(t, gock.IsDone())
	})

	t.Run("WithValidResponse", func(t *testing.T) {
		defer gock.Off() // Flush pending mocks after test execution

		response, err := os.ReadFile("testdata/mark_as_read_200.txt")
		require.NoError(t, err)

		// Mock the request
		gock.New("https://freshrss.example.com").
			Post("/reader/api/0/mark-all-as-read").
			MatchHeader("Authorization", "GoogleLogin auth=test/auth-token").
			MatchHeader("Content-Type", "application/x-www-form-urlencoded").
			Reply(200).
			BodyString(string(response))

		c := initTestClient(t)

		err = c.MarkAsRead(context.Background(), "test/auth-token", "feed-id", 7)
		require.NoError(t, err)
		assert.True(t, gock.IsDone())
	})
}
