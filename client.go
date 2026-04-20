// Package elevenlabs is a tiny Go client for the ElevenLabs Agents API.
//
// It covers the minimum surface needed to test a voice agent end-to-end over
// telephony: create/update an agent, import a Twilio phone number, place an
// outbound call, fetch conversation results, and mint a signed URL for
// browser/native testing.
package elevenlabs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultBaseURL = "https://api.elevenlabs.io"

// Client talks to the ElevenLabs REST API.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// Option customises a Client.
type Option func(*Client)

// WithBaseURL overrides the API base URL (e.g. a regional residency endpoint).
func WithBaseURL(url string) Option {
	return func(c *Client) {
		if strings.TrimSpace(url) != "" {
			c.baseURL = strings.TrimRight(url, "/")
		}
	}
}

// WithHTTPClient swaps the underlying *http.Client (timeouts, transport, etc).
func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) {
		if h != nil {
			c.httpClient = h
		}
	}
}

// New constructs a Client. apiKey is your ElevenLabs xi-api-key.
func New(apiKey string, opts ...Option) (*Client, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("elevenlabs: api key is required")
	}
	c := &Client{
		apiKey:     apiKey,
		baseURL:    defaultBaseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

// APIError is returned for non-2xx responses.
type APIError struct {
	StatusCode int
	Body       string
}

// Error implements the error interface.
func (e *APIError) Error() string {
	return fmt.Sprintf("elevenlabs: status %d: %s", e.StatusCode, strings.TrimSpace(e.Body))
}

// do issues an authenticated JSON request and decodes the response into out
// (which may be nil to discard the body).
func (c *Client) do(ctx context.Context, method, path string, body, out any) error {
	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("elevenlabs: marshal request: %w", err)
		}
		reader = bytes.NewReader(raw)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reader)
	if err != nil {
		return fmt.Errorf("elevenlabs: build request: %w", err)
	}
	req.Header.Set("xi-api-key", c.apiKey)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("elevenlabs: request: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("elevenlabs: read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return &APIError{StatusCode: resp.StatusCode, Body: string(raw)}
	}

	if out == nil || len(raw) == 0 {
		return nil
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("elevenlabs: decode response: %w", err)
	}
	return nil
}
