package serper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"mbvlabs/config"
)

const (
	googleURL = "https://google.serper.dev"
	scrapeURL = "https://scrape.serper.dev"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
	googleURL  string
	scrapeURL  string
}

func New(cfg config.Config) *Client {
	return &Client{
		apiKey:     cfg.Serper.APIKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		googleURL:  googleURL,
		scrapeURL:  scrapeURL,
	}
}

type Query struct {
	Query        string `json:"q"`
	Location     string `json:"location,omitempty"`
	CountryCode  string `json:"gl,omitempty"`
	LanguageCode string `json:"hl,omitempty"`
	Autocorrect  *bool  `json:"autocorrect,omitempty"`
	Page         int    `json:"page,omitempty"`
	Num          int    `json:"num,omitempty"`
	TimeRange    string `json:"tbs,omitempty"`
}

type SearchParameters struct {
	Query  string `json:"q"`
	Type   string `json:"type"`
	Engine string `json:"engine"`
}

type SearchResult struct {
	Title     string     `json:"title"`
	Link      string     `json:"link"`
	Snippet   string     `json:"snippet"`
	Date      string     `json:"date"`
	Position  int        `json:"position"`
	Sitelinks []Sitelink `json:"sitelinks,omitempty"`
}

type Sitelink struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

type Question struct {
	Question string `json:"question"`
	Snippet  string `json:"snippet"`
	Title    string `json:"title"`
	Link     string `json:"link"`
}

type RelatedSearch struct {
	Query string `json:"query"`
}

type SearchResponse struct {
	SearchParameters SearchParameters `json:"searchParameters"`
	KnowledgeGraph   map[string]any   `json:"knowledgeGraph,omitempty"`
	Organic          []SearchResult   `json:"organic"`
	PeopleAlsoAsk    []Question       `json:"peopleAlsoAsk,omitempty"`
	RelatedSearches  []RelatedSearch  `json:"relatedSearches,omitempty"`
	Credits          int              `json:"credits"`
}

type NewsResult struct {
	Title    string `json:"title"`
	Link     string `json:"link"`
	Snippet  string `json:"snippet"`
	Date     string `json:"date"`
	Source   string `json:"source"`
	ImageURL string `json:"imageUrl"`
	Position int    `json:"position"`
}

type NewsResponse struct {
	SearchParameters SearchParameters `json:"searchParameters"`
	News             []NewsResult     `json:"news"`
	Credits          int              `json:"credits"`
}

type ScrapeRequest struct {
	URL             string `json:"url"`
	IncludeMarkdown bool   `json:"includeMarkdown,omitempty"`
}

type ScrapeResponse struct {
	Text     string          `json:"text"`
	Markdown string          `json:"markdown,omitempty"`
	Metadata map[string]any  `json:"metadata,omitempty"`
	JSONLD   json.RawMessage `json:"jsonld,omitempty"`
}

func (c *Client) Search(ctx context.Context, query Query) (SearchResponse, error) {
	var response SearchResponse
	if strings.TrimSpace(query.Query) == "" {
		return response, fmt.Errorf("serper: search query is required")
	}
	return response, c.post(ctx, c.googleURL+"/search", query, &response)
}

func (c *Client) News(ctx context.Context, query Query) (NewsResponse, error) {
	var response NewsResponse
	if strings.TrimSpace(query.Query) == "" {
		return response, fmt.Errorf("serper: news query is required")
	}
	return response, c.post(ctx, c.googleURL+"/news", query, &response)
}

func (c *Client) Scrape(ctx context.Context, request ScrapeRequest) (ScrapeResponse, error) {
	var response ScrapeResponse
	if strings.TrimSpace(request.URL) == "" {
		return response, fmt.Errorf("serper: scrape URL is required")
	}
	return response, c.post(ctx, c.scrapeURL, request, &response)
}

func (c *Client) post(ctx context.Context, url string, payload, response any) error {
	if strings.TrimSpace(c.apiKey) == "" {
		return fmt.Errorf("serper: API key is required")
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("serper: encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("serper: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", c.apiKey)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("serper: send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		message, _ := io.ReadAll(io.LimitReader(res.Body, 64<<10))
		return fmt.Errorf("serper: %s: %s", res.Status, strings.TrimSpace(string(message)))
	}

	if err := json.NewDecoder(res.Body).Decode(response); err != nil {
		return fmt.Errorf("serper: decode response: %w", err)
	}
	return nil
}
