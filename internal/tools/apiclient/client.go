package apiclient

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
)

type Client struct {
	httpClient *http.Client

	mu        sync.Mutex
	cachePath string
	cache     map[string][]byte
}

type Options func(*Client) error

func WithTimeout(timeout time.Duration) Options {
	return func(c *Client) error {
		c.httpClient.Timeout = timeout
		return nil
	}
}

func WithCache(filePath string) Options {
	return func(c *Client) error {
		c.cachePath = filePath
		return c.loadCache()
	}
}

func New(opts ...Options) (*Client, error) {
	c := &Client{
		httpClient: &http.Client{},
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return c, nil
}

func (c *Client) Get(
	ctx context.Context,
	reqURL string,
	headers map[string]string,
	res any,
) error {
	if c.cache[reqURL] != nil {
		return json.Unmarshal(c.cache[reqURL], res)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	err = c.UpdateCache(reqURL, bs)
	if err != nil {
		return fmt.Errorf("failed to update cache: %w", err)
	}

	return json.Unmarshal(bs, res)
}

// TODO: maybe bucket storage or redis?

func (c *Client) UpdateCache(key string, value []byte) error {
	if c.cachePath == "" {
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cache == nil {
		c.cache = make(map[string][]byte)
	}

	c.cache[key] = value

	f, err := os.OpenFile(c.cachePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open cache file: %w", err)
	}
	defer f.Close()

	if err := gob.NewEncoder(f).Encode(c.cache); err != nil {
		return fmt.Errorf("failed to encode cache: %w", err)
	}

	return nil
}

func (c *Client) loadCache() error {
	f, err := os.OpenFile(c.cachePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open cache file: %w", err)
	}
	defer f.Close()

	cache := make(map[string][]byte)

	err = gob.NewDecoder(f).Decode(&cache)
	if err != nil {
		// TODO: maybe try a different approach
		c.cache = make(map[string][]byte) // if decoding fails, return an empty cache
		return nil
	}
	c.cache = cache

	slog.Info("Cache loaded", "size", len(cache))

	return nil
}
