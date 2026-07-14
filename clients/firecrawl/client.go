package firecrawl

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

	"mbvlabs/config"
)

const apiURL = "https://api.firecrawl.dev/v2"

type Client struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

func New(cfg config.Config) *Client {
	return &Client{
		apiKey:     cfg.Firecrawl.APIKey,
		httpClient: &http.Client{Timeout: 5 * time.Minute},
		baseURL:    apiURL,
	}
}

type Location struct {
	Country   string   `json:"country,omitempty"`
	Languages []string `json:"languages,omitempty"`
}

type Parser struct {
	Type     string `json:"type"`
	Mode     string `json:"mode,omitempty"`
	MaxPages int    `json:"maxPages,omitempty"`
}

type Profile struct {
	Name        string `json:"name"`
	SaveChanges *bool  `json:"saveChanges,omitempty"`
}

type ScrapeOptions struct {
	Formats             []any             `json:"formats,omitempty"`
	OnlyMainContent     *bool             `json:"onlyMainContent,omitempty"`
	OnlyCleanContent    bool              `json:"onlyCleanContent,omitempty"`
	IncludeTags         []string          `json:"includeTags,omitempty"`
	ExcludeTags         []string          `json:"excludeTags,omitempty"`
	MaxAge              *int              `json:"maxAge,omitempty"`
	MinAge              *int              `json:"minAge,omitempty"`
	Headers             map[string]string `json:"headers,omitempty"`
	WaitFor             int               `json:"waitFor,omitempty"`
	Mobile              bool              `json:"mobile,omitempty"`
	SkipTLSVerification *bool             `json:"skipTlsVerification,omitempty"`
	Timeout             int               `json:"timeout,omitempty"`
	Parsers             []Parser          `json:"parsers,omitempty"`
	Actions             []map[string]any  `json:"actions,omitempty"`
	Location            *Location         `json:"location,omitempty"`
	RemoveBase64Images  *bool             `json:"removeBase64Images,omitempty"`
	BlockAds            *bool             `json:"blockAds,omitempty"`
	Proxy               string            `json:"proxy,omitempty"`
	StoreInCache        *bool             `json:"storeInCache,omitempty"`
	Lockdown            bool              `json:"lockdown,omitempty"`
	RedactPII           any               `json:"redactPII,omitempty"`
	Profile             *Profile          `json:"profile,omitempty"`
	ThreatProtection    map[string]any    `json:"threatProtection,omitempty"`
}

type ScrapeRequest struct {
	URL string `json:"url"`
	ScrapeOptions
	ZeroDataRetention bool `json:"zeroDataRetention,omitempty"`
}

type Document struct {
	Title       string          `json:"title,omitempty"`
	Description string          `json:"description,omitempty"`
	URL         string          `json:"url,omitempty"`
	Markdown    string          `json:"markdown,omitempty"`
	Summary     string          `json:"summary,omitempty"`
	HTML        string          `json:"html,omitempty"`
	RawHTML     string          `json:"rawHtml,omitempty"`
	Screenshot  string          `json:"screenshot,omitempty"`
	Audio       string          `json:"audio,omitempty"`
	Video       string          `json:"video,omitempty"`
	Answer      string          `json:"answer,omitempty"`
	Highlights  string          `json:"highlights,omitempty"`
	Links       []string        `json:"links,omitempty"`
	JSON        json.RawMessage `json:"json,omitempty"`
	Metadata    map[string]any  `json:"metadata,omitempty"`
	Actions     json.RawMessage `json:"actions,omitempty"`
	Branding    json.RawMessage `json:"branding,omitempty"`
	Product     json.RawMessage `json:"product,omitempty"`
	Menu        json.RawMessage `json:"menu,omitempty"`
	Warning     string          `json:"warning,omitempty"`
}

type ScrapeResponse struct {
	Success bool     `json:"success"`
	Data    Document `json:"data"`
}

type Webhook struct {
	URL      string            `json:"url"`
	Headers  map[string]string `json:"headers,omitempty"`
	Metadata map[string]any    `json:"metadata,omitempty"`
	Events   []string          `json:"events,omitempty"`
}

type CrawlRequest struct {
	URL                   string         `json:"url"`
	Prompt                string         `json:"prompt,omitempty"`
	ExcludePaths          []string       `json:"excludePaths,omitempty"`
	IncludePaths          []string       `json:"includePaths,omitempty"`
	MaxDiscoveryDepth     *int           `json:"maxDiscoveryDepth,omitempty"`
	Sitemap               string         `json:"sitemap,omitempty"`
	IgnoreQueryParameters bool           `json:"ignoreQueryParameters,omitempty"`
	RegexOnFullURL        bool           `json:"regexOnFullURL,omitempty"`
	Limit                 int            `json:"limit,omitempty"`
	CrawlEntireDomain     bool           `json:"crawlEntireDomain,omitempty"`
	AllowExternalLinks    bool           `json:"allowExternalLinks,omitempty"`
	AllowSubdomains       bool           `json:"allowSubdomains,omitempty"`
	IgnoreRobotsTxt       bool           `json:"ignoreRobotsTxt,omitempty"`
	RobotsUserAgent       string         `json:"robotsUserAgent,omitempty"`
	Delay                 float64        `json:"delay,omitempty"`
	MaxConcurrency        int            `json:"maxConcurrency,omitempty"`
	Webhook               *Webhook       `json:"webhook,omitempty"`
	ScrapeOptions         *ScrapeOptions `json:"scrapeOptions,omitempty"`
	ZeroDataRetention     bool           `json:"zeroDataRetention,omitempty"`
}

type CrawlResponse struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
	URL     string `json:"url"`
}

type CrawlStatusResponse struct {
	Status      string     `json:"status"`
	Total       int        `json:"total"`
	Completed   int        `json:"completed"`
	CreditsUsed int        `json:"creditsUsed"`
	ExpiresAt   string     `json:"expiresAt"`
	CreatedAt   string     `json:"createdAt"`
	CompletedAt string     `json:"completedAt,omitempty"`
	Duration    float64    `json:"duration"`
	Next        string     `json:"next,omitempty"`
	Data        []Document `json:"data"`
}

type MapRequest struct {
	URL                   string         `json:"url"`
	Search                string         `json:"search,omitempty"`
	Sitemap               string         `json:"sitemap,omitempty"`
	IncludeSubdomains     *bool          `json:"includeSubdomains,omitempty"`
	IgnoreQueryParameters *bool          `json:"ignoreQueryParameters,omitempty"`
	IgnoreCache           bool           `json:"ignoreCache,omitempty"`
	Limit                 int            `json:"limit,omitempty"`
	Timeout               int            `json:"timeout,omitempty"`
	Location              *Location      `json:"location,omitempty"`
	ThreatProtection      map[string]any `json:"threatProtection,omitempty"`
}

type MapLink struct {
	URL         string `json:"url"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type MapResponse struct {
	Success bool      `json:"success"`
	Links   []MapLink `json:"links"`
}

type SearchSource struct {
	Type     string `json:"type"`
	TBS      string `json:"tbs,omitempty"`
	Location string `json:"location,omitempty"`
}

type SearchCategory struct {
	Type string `json:"type"`
}

type SearchRequest struct {
	Query             string           `json:"query"`
	Limit             int              `json:"limit,omitempty"`
	Sources           []SearchSource   `json:"sources,omitempty"`
	Categories        []SearchCategory `json:"categories,omitempty"`
	IncludeDomains    []string         `json:"includeDomains,omitempty"`
	ExcludeDomains    []string         `json:"excludeDomains,omitempty"`
	TBS               string           `json:"tbs,omitempty"`
	Location          string           `json:"location,omitempty"`
	Country           string           `json:"country,omitempty"`
	Timeout           int              `json:"timeout,omitempty"`
	IgnoreInvalidURLs bool             `json:"ignoreInvalidURLs,omitempty"`
	Enterprise        []string         `json:"enterprise,omitempty"`
	ScrapeOptions     *ScrapeOptions   `json:"scrapeOptions,omitempty"`
	ThreatProtection  map[string]any   `json:"threatProtection,omitempty"`
}

type ImageResult struct {
	Title       string `json:"title"`
	ImageURL    string `json:"imageUrl"`
	ImageWidth  int    `json:"imageWidth"`
	ImageHeight int    `json:"imageHeight"`
	URL         string `json:"url"`
	Position    int    `json:"position"`
}

type NewsResult struct {
	Document
	Snippet  string `json:"snippet,omitempty"`
	Date     string `json:"date,omitempty"`
	ImageURL string `json:"imageUrl,omitempty"`
	Position int    `json:"position,omitempty"`
}

type SearchResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Web    []Document    `json:"web,omitempty"`
		Images []ImageResult `json:"images,omitempty"`
		News   []NewsResult  `json:"news,omitempty"`
	} `json:"data"`
	Warning     string `json:"warning,omitempty"`
	ID          string `json:"id"`
	CreditsUsed int    `json:"creditsUsed"`
}

func (c *Client) Scrape(ctx context.Context, request ScrapeRequest) (ScrapeResponse, error) {
	var response ScrapeResponse
	if strings.TrimSpace(request.URL) == "" {
		return response, fmt.Errorf("firecrawl: scrape URL is required")
	}
	return response, c.do(ctx, http.MethodPost, "/scrape", request, &response)
}

func (c *Client) Crawl(ctx context.Context, request CrawlRequest) (CrawlResponse, error) {
	var response CrawlResponse
	if strings.TrimSpace(request.URL) == "" {
		return response, fmt.Errorf("firecrawl: crawl URL is required")
	}
	return response, c.do(ctx, http.MethodPost, "/crawl", request, &response)
}

func (c *Client) CrawlStatus(ctx context.Context, id string) (CrawlStatusResponse, error) {
	var response CrawlStatusResponse
	if strings.TrimSpace(id) == "" {
		return response, fmt.Errorf("firecrawl: crawl ID is required")
	}
	return response, c.do(ctx, http.MethodGet, "/crawl/"+url.PathEscape(id), nil, &response)
}

func (c *Client) Map(ctx context.Context, request MapRequest) (MapResponse, error) {
	var response MapResponse
	if strings.TrimSpace(request.URL) == "" {
		return response, fmt.Errorf("firecrawl: map URL is required")
	}
	return response, c.do(ctx, http.MethodPost, "/map", request, &response)
}

func (c *Client) Search(ctx context.Context, request SearchRequest) (SearchResponse, error) {
	var response SearchResponse
	if strings.TrimSpace(request.Query) == "" {
		return response, fmt.Errorf("firecrawl: search query is required")
	}
	return response, c.do(ctx, http.MethodPost, "/search", request, &response)
}

func (c *Client) do(ctx context.Context, method, path string, payload, response any) error {
	if strings.TrimSpace(c.apiKey) == "" {
		return fmt.Errorf("firecrawl: API key is required")
	}

	var body io.Reader
	if payload != nil {
		encoded, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("firecrawl: encode request: %w", err)
		}
		body = bytes.NewReader(encoded)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return fmt.Errorf("firecrawl: create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("firecrawl: send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		message, _ := io.ReadAll(io.LimitReader(res.Body, 64<<10))
		return fmt.Errorf("firecrawl: %s: %s", res.Status, strings.TrimSpace(string(message)))
	}

	if err := json.NewDecoder(res.Body).Decode(response); err != nil {
		return fmt.Errorf("firecrawl: decode response: %w", err)
	}
	return nil
}
