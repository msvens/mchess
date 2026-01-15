package upstream

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// Client handles requests to the upstream schack.se API
type Client struct {
	baseURL    string
	httpClient *http.Client
	limiter    *rate.Limiter
}

// NewClient creates a new upstream API client
func NewClient(baseURL string, timeout time.Duration, rateLimit int) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		limiter: rate.NewLimiter(rate.Limit(rateLimit), rateLimit), // allow bursts up to rateLimit
	}
}

// get performs a rate-limited GET request and decodes the response
func (c *Client) get(ctx context.Context, path string, result interface{}) error {
	body, err := c.GetRaw(ctx, path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

// GetRaw performs a rate-limited GET request and returns raw bytes (for pass-through)
func (c *Client) GetRaw(ctx context.Context, path string) ([]byte, error) {
	// Wait for rate limiter
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	url := c.baseURL + path
	slog.Debug("Upstream request", "url", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upstream error: status=%d body=%s", resp.StatusCode, string(body))
	}

	return body, nil
}
