package services

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebSearchClientProviderRequests(t *testing.T) {
	var method, path, auth, key string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, path, auth, key = r.Method, r.URL.EscapedPath(), r.Header.Get("Authorization"), r.Header.Get("X-API-KEY")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"id":"123"}`))
	}))
	defer server.Close()
	client := NewWebSearchClient("serper-key", "firecrawl-key")
	client.SerperURL, client.FirecrawlURL = server.URL, server.URL

	result, err := client.Request(context.Background(), "search", "", map[string]string{"q": "test"})
	if err != nil || !json.Valid(result) || method != "POST" || path != "/search" || key != "serper-key" || auth != "" {
		t.Fatalf("search = %s, %v, %s %s auth=%q key=%q", result, err, method, path, auth, key)
	}
	result, err = client.Request(context.Background(), "crawl status", "one/two", nil)
	if err != nil || !json.Valid(result) || method != "GET" || path != "/crawl/one%2Ftwo" || auth != "Bearer firecrawl-key" {
		t.Fatalf("crawl = %s, %v, %s %s auth=%q", result, err, method, path, auth)
	}
	client.FirecrawlKey = ""
	if _, err := client.Request(context.Background(), "scrape", "", nil); err == nil {
		t.Fatal("empty API key accepted")
	}
}
