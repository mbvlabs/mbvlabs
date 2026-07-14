package serper

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.Header.Get("X-API-KEY") != "secret" {
			t.Errorf("X-API-KEY = %q, want secret", r.Header.Get("X-API-KEY"))
		}

		switch r.URL.Path {
		case "/search":
			json.NewEncoder(w).Encode(SearchResponse{Organic: []SearchResult{{Title: "Result"}}})
		case "/news":
			json.NewEncoder(w).Encode(NewsResponse{News: []NewsResult{{Title: "Headline"}}})
		case "/scrape":
			json.NewEncoder(w).Encode(ScrapeResponse{Text: "Page", Markdown: "# Page"})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := &Client{
		apiKey:     "secret",
		httpClient: server.Client(),
		googleURL:  server.URL,
		scrapeURL:  server.URL + "/scrape",
	}

	search, err := client.Search(context.Background(), Query{Query: "golang"})
	if err != nil || len(search.Organic) != 1 || search.Organic[0].Title != "Result" {
		t.Fatalf("Search() = %#v, %v", search, err)
	}

	news, err := client.News(context.Background(), Query{Query: "golang"})
	if err != nil || len(news.News) != 1 || news.News[0].Title != "Headline" {
		t.Fatalf("News() = %#v, %v", news, err)
	}

	scrape, err := client.Scrape(context.Background(), ScrapeRequest{
		URL:             "https://example.com",
		IncludeMarkdown: true,
	})
	if err != nil || scrape.Markdown != "# Page" {
		t.Fatalf("Scrape() = %#v, %v", scrape, err)
	}
}
