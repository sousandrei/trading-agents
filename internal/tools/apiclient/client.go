package apiclient

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type Client struct {
	httpClient *http.Client

	mu    sync.Mutex
	cache map[string][]byte
}

func New(
	timeout time.Duration,
) (*Client, error) {
	cache, err := loadCache()
	if err != nil {
		return nil, fmt.Errorf("failed to load cache: %w", err)
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		cache: cache,
	}, nil
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

// TODO: toggable cache mechanism by environment variable, maybe bucket storage or redis?

const cachePath = "data/cache.bin"

func (c *Client) UpdateCache(key string, value []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cache == nil {
		c.cache = make(map[string][]byte)
	}

	c.cache[key] = value

	f, err := os.OpenFile(cachePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open cache file: %w", err)
	}
	defer f.Close()

	if err := gob.NewEncoder(f).Encode(c.cache); err != nil {
		return fmt.Errorf("failed to encode cache: %w", err)
	}

	return nil
}

func loadCache() (map[string][]byte, error) {
	f, err := os.OpenFile(cachePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open cache file: %w", err)
	}
	defer f.Close()

	cache := make(map[string][]byte)

	err = gob.NewDecoder(f).Decode(&cache)
	if err != nil {
		return nil, fmt.Errorf("failed to decode cache: %w", err)
	}

	return cache, nil
}
