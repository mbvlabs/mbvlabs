package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"mbvlabs/routes"
	"mbvlabs/services"

	"github.com/mbvlabs/andurel/pkg/anduin"
)

func TestAPISearchRoutes(t *testing.T) {
	var upstream map[string]any
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if err := json.NewDecoder(r.Body).Decode(&upstream); err != nil {
				t.Error(err)
			}
		}
		_, _ = w.Write([]byte(`{"success":true,"id":"123"}`))
	}))
	defer provider.Close()
	client := services.NewWebSearchClient("s", "f")
	client.SerperURL, client.FirecrawlURL = provider.URL, provider.URL
	rt := anduin.NewRouter()
	NewAPISearchHandler(client, "user", "pass").Register(rt)
	request := func(method, path, body string, auth bool) *httptest.ResponseRecorder {
		req := httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		if auth {
			req.SetBasicAuth("user", "pass")
		}
		rec := httptest.NewRecorder()
		rt.Handler().ServeHTTP(rec, req)
		return rec
	}
	if got := request("POST", routes.ApiSearch.Path(), `{"query":"q"}`, false); got.Code != 401 {
		t.Fatalf("unauthorized status = %d", got.Code)
	}
	if got := request("POST", routes.ApiSearch.Path(), `{"query":" "}`, true); got.Code != 400 || !strings.Contains(got.Body.String(), "query is required") {
		t.Fatalf("missing query = %d %s", got.Code, got.Body)
	}
	if got := request("POST", routes.ApiSearch.Path(), `{"query":"term","countryCode":"es"}`, true); got.Code != 200 || upstream["q"] != "term" || upstream["gl"] != "es" {
		t.Fatalf("search = %d %s, upstream=%v", got.Code, got.Body, upstream)
	}
	if got := request("POST", routes.ApiCrawl.Path(), `{"url":"https://example.com"}`, true); got.Code != 202 {
		t.Fatalf("crawl = %d %s", got.Code, got.Body)
	}
	if got := request("GET", "/api/crawl/123", "", true); got.Code != 200 {
		t.Fatalf("crawl status = %d %s", got.Code, got.Body)
	}
	if got := request("POST", routes.ApiScrape.Path(), `{"url":"https://example.com","waitFor":"bad"}`, true); got.Code != 400 {
		t.Fatalf("invalid scrape = %d %s", got.Code, got.Body)
	}
}
