package firecrawl

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer secret" {
			t.Errorf("Authorization = %q, want Bearer secret", r.Header.Get("Authorization"))
		}

		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/scrape":
			json.NewEncoder(w).Encode(ScrapeResponse{Success: true, Data: Document{Markdown: "# Page"}})
		case r.Method == http.MethodPost && r.URL.Path == "/crawl":
			json.NewEncoder(w).Encode(CrawlResponse{Success: true, ID: "crawl-id"})
		case r.Method == http.MethodGet && r.URL.Path == "/crawl/crawl-id":
			json.NewEncoder(w).Encode(CrawlStatusResponse{Status: "completed", Completed: 1})
		case r.Method == http.MethodPost && r.URL.Path == "/map":
			json.NewEncoder(w).Encode(MapResponse{Success: true, Links: []MapLink{{URL: "https://example.com"}}})
		case r.Method == http.MethodPost && r.URL.Path == "/search":
			json.NewEncoder(w).Encode(map[string]any{
				"success": true,
				"data":    map[string]any{"web": []map[string]any{{"title": "Result"}}},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := &Client{apiKey: "secret", httpClient: server.Client(), baseURL: server.URL}
	ctx := context.Background()

	scrape, err := client.Scrape(ctx, ScrapeRequest{URL: "https://example.com"})
	if err != nil || scrape.Data.Markdown != "# Page" {
		t.Fatalf("Scrape() = %#v, %v", scrape, err)
	}

	crawl, err := client.Crawl(ctx, CrawlRequest{URL: "https://example.com"})
	if err != nil || crawl.ID != "crawl-id" {
		t.Fatalf("Crawl() = %#v, %v", crawl, err)
	}

	status, err := client.CrawlStatus(ctx, crawl.ID)
	if err != nil || status.Status != "completed" {
		t.Fatalf("CrawlStatus() = %#v, %v", status, err)
	}

	mapped, err := client.Map(ctx, MapRequest{URL: "https://example.com"})
	if err != nil || len(mapped.Links) != 1 {
		t.Fatalf("Map() = %#v, %v", mapped, err)
	}

	search, err := client.Search(ctx, SearchRequest{Query: "golang"})
	if err != nil || len(search.Data.Web) != 1 || search.Data.Web[0].Title != "Result" {
		t.Fatalf("Search() = %#v, %v", search, err)
	}
}
