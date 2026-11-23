package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// WebSearchClient calls Serper for search and Firecrawl for scrape, crawl and map.
type WebSearchClient struct {
	SerperKey, FirecrawlKey string
	SerperURL, FirecrawlURL string
	HTTPClient              *http.Client
}

func NewWebSearchClient(serperKey, firecrawlKey string) *WebSearchClient {
	return &WebSearchClient{
		SerperKey: serperKey, FirecrawlKey: firecrawlKey,
		SerperURL: "https://google.serper.dev", FirecrawlURL: "https://api.firecrawl.dev/v2",
		HTTPClient: &http.Client{Timeout: 5 * time.Minute},
	}
}

func (c *WebSearchClient) Request(ctx context.Context, operation, id string, payload any) (json.RawMessage, error) {
	base, key, method, path := c.FirecrawlURL, c.FirecrawlKey, http.MethodPost, "/"+operation
	if operation == "search" {
		base, key = c.SerperURL, c.SerperKey
	}
	if operation == "crawl status" {
		method, path = http.MethodGet, "/crawl/"+url.PathEscape(id)
	}
	if strings.TrimSpace(key) == "" {
		return nil, fmt.Errorf("%s: API key is required", operation)
	}
	if operation == "search" {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
	}
	var body io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, method, base+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if operation == "search" {
		req.Header.Set("X-API-KEY", key)
	} else {
		req.Header.Set("Authorization", "Bearer "+key)
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("%s provider returned %s", operation, res.Status)
	}
	var result json.RawMessage
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}
