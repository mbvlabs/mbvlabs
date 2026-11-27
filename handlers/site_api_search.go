package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"mbvlabs/middleware"
	"mbvlabs/routes"
	"mbvlabs/services"

	"github.com/mbvlabs/andurel/pkg/anduin"
)

type APISearchHandler struct {
	client             *services.WebSearchClient
	username, password string
}

func NewAPISearchHandler(client *services.WebSearchClient, username, password string) APISearchHandler {
	return APISearchHandler{client: client, username: username, password: password}
}

func (h APISearchHandler) Register(rt *anduin.Router) {
	auth := middleware.APIBasicAuth(h.username, h.password)
	rt.Route(routes.ApiSearch).Post(h.search, auth)
	rt.Route(routes.ApiScrape).Post(h.post("scrape"), auth)
	rt.Route(routes.ApiCrawl).Post(h.post("crawl"), auth)
	rt.Route(routes.ApiCrawlStatus).Get(h.crawlStatus, auth)
	rt.Route(routes.ApiMap).Post(h.post("map"), auth)
}

type apiError struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

type webSearchRequest struct {
	Query        string `json:"query"`
	Location     string `json:"location,omitempty"`
	CountryCode  string `json:"countryCode,omitempty"`
	LanguageCode string `json:"languageCode,omitempty"`
	Autocorrect  *bool  `json:"autocorrect,omitempty"`
	Page         int    `json:"page,omitempty"`
	Num          int    `json:"num,omitempty"`
	TimeRange    string `json:"timeRange,omitempty"`
}

type apiLocation struct {
	Country   string   `json:"country,omitempty"`
	Languages []string `json:"languages,omitempty"`
}
type apiParser struct {
	Type     string `json:"type"`
	Mode     string `json:"mode,omitempty"`
	MaxPages int    `json:"maxPages,omitempty"`
}
type apiProfile struct {
	Name        string `json:"name"`
	SaveChanges *bool  `json:"saveChanges,omitempty"`
}
type apiScrapeOptions struct {
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
	Parsers             []apiParser       `json:"parsers,omitempty"`
	Actions             []map[string]any  `json:"actions,omitempty"`
	Location            *apiLocation      `json:"location,omitempty"`
	RemoveBase64Images  *bool             `json:"removeBase64Images,omitempty"`
	BlockAds            *bool             `json:"blockAds,omitempty"`
	Proxy               string            `json:"proxy,omitempty"`
	StoreInCache        *bool             `json:"storeInCache,omitempty"`
	Lockdown            bool              `json:"lockdown,omitempty"`
	RedactPII           any               `json:"redactPII,omitempty"`
	Profile             *apiProfile       `json:"profile,omitempty"`
	ThreatProtection    map[string]any    `json:"threatProtection,omitempty"`
}
type apiScrapeRequest struct {
	URL string `json:"url"`
	apiScrapeOptions
	ZeroDataRetention bool `json:"zeroDataRetention,omitempty"`
}
type apiWebhook struct {
	URL      string            `json:"url"`
	Headers  map[string]string `json:"headers,omitempty"`
	Metadata map[string]any    `json:"metadata,omitempty"`
	Events   []string          `json:"events,omitempty"`
}
type apiCrawlRequest struct {
	URL                   string            `json:"url"`
	Prompt                string            `json:"prompt,omitempty"`
	ExcludePaths          []string          `json:"excludePaths,omitempty"`
	IncludePaths          []string          `json:"includePaths,omitempty"`
	MaxDiscoveryDepth     *int              `json:"maxDiscoveryDepth,omitempty"`
	Sitemap               string            `json:"sitemap,omitempty"`
	IgnoreQueryParameters bool              `json:"ignoreQueryParameters,omitempty"`
	RegexOnFullURL        bool              `json:"regexOnFullURL,omitempty"`
	Limit                 int               `json:"limit,omitempty"`
	CrawlEntireDomain     bool              `json:"crawlEntireDomain,omitempty"`
	AllowExternalLinks    bool              `json:"allowExternalLinks,omitempty"`
	AllowSubdomains       bool              `json:"allowSubdomains,omitempty"`
	IgnoreRobotsTxt       bool              `json:"ignoreRobotsTxt,omitempty"`
	RobotsUserAgent       string            `json:"robotsUserAgent,omitempty"`
	Delay                 float64           `json:"delay,omitempty"`
	MaxConcurrency        int               `json:"maxConcurrency,omitempty"`
	Webhook               *apiWebhook       `json:"webhook,omitempty"`
	ScrapeOptions         *apiScrapeOptions `json:"scrapeOptions,omitempty"`
	ZeroDataRetention     bool              `json:"zeroDataRetention,omitempty"`
}
type apiMapRequest struct {
	URL                   string         `json:"url"`
	Search                string         `json:"search,omitempty"`
	Sitemap               string         `json:"sitemap,omitempty"`
	IncludeSubdomains     *bool          `json:"includeSubdomains,omitempty"`
	IgnoreQueryParameters *bool          `json:"ignoreQueryParameters,omitempty"`
	IgnoreCache           bool           `json:"ignoreCache,omitempty"`
	Limit                 int            `json:"limit,omitempty"`
	Timeout               int            `json:"timeout,omitempty"`
	Location              *apiLocation   `json:"location,omitempty"`
	ThreatProtection      map[string]any `json:"threatProtection,omitempty"`
}

type apiDocument struct {
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
type apiSearchResponse struct {
	SearchParameters struct {
		Query  string `json:"q"`
		Type   string `json:"type"`
		Engine string `json:"engine"`
	} `json:"searchParameters"`
	KnowledgeGraph map[string]any `json:"knowledgeGraph,omitempty"`
	Organic        []struct {
		Title     string `json:"title"`
		Link      string `json:"link"`
		Snippet   string `json:"snippet"`
		Date      string `json:"date"`
		Position  int    `json:"position"`
		Sitelinks []struct {
			Title string `json:"title"`
			Link  string `json:"link"`
		} `json:"sitelinks,omitempty"`
	} `json:"organic"`
	PeopleAlsoAsk []struct {
		Question string `json:"question"`
		Snippet  string `json:"snippet"`
		Title    string `json:"title"`
		Link     string `json:"link"`
	} `json:"peopleAlsoAsk,omitempty"`
	RelatedSearches []struct {
		Query string `json:"query"`
	} `json:"relatedSearches,omitempty"`
	Credits int `json:"credits"`
}
type apiScrapeResponse struct {
	Success bool        `json:"success"`
	Data    apiDocument `json:"data"`
}
type apiCrawlResponse struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
	URL     string `json:"url"`
}
type apiCrawlStatusResponse struct {
	Status      string        `json:"status"`
	Total       int           `json:"total"`
	Completed   int           `json:"completed"`
	CreditsUsed int           `json:"creditsUsed"`
	ExpiresAt   string        `json:"expiresAt"`
	CreatedAt   string        `json:"createdAt"`
	CompletedAt string        `json:"completedAt,omitempty"`
	Duration    float64       `json:"duration"`
	Next        string        `json:"next,omitempty"`
	Data        []apiDocument `json:"data"`
}
type apiMapResponse struct {
	Success bool `json:"success"`
	Links   []struct {
		URL         string `json:"url"`
		Title       string `json:"title,omitempty"`
		Description string `json:"description,omitempty"`
	} `json:"links"`
}

func (h APISearchHandler) search(ctx context.Context) error {
	var p webSearchRequest
	if err := anduin.DecodeJSON(ctx, &p); err != nil {
		return apiFailure(ctx, 400, "invalid JSON request body")
	}
	if strings.TrimSpace(p.Query) == "" {
		return apiFailure(ctx, 400, "query is required")
	}
	request := struct {
		Query        string `json:"q"`
		Location     string `json:"location,omitempty"`
		CountryCode  string `json:"gl,omitempty"`
		LanguageCode string `json:"hl,omitempty"`
		Autocorrect  *bool  `json:"autocorrect,omitempty"`
		Page         int    `json:"page,omitempty"`
		Num          int    `json:"num,omitempty"`
		TimeRange    string `json:"tbs,omitempty"`
	}{p.Query, p.Location, p.CountryCode, p.LanguageCode, p.Autocorrect, p.Page, p.Num, p.TimeRange}
	return h.call(ctx, "search", "", request)
}

func (h APISearchHandler) post(operation string) anduin.Handler {
	return func(ctx context.Context) error {
		var payload any
		switch operation {
		case "scrape":
			payload = &apiScrapeRequest{}
		case "crawl":
			payload = &apiCrawlRequest{}
		case "map":
			payload = &apiMapRequest{}
		}
		if err := anduin.DecodeJSON(ctx, payload); err != nil {
			return apiFailure(ctx, 400, "invalid JSON request body")
		}
		var url string
		switch p := payload.(type) {
		case *apiScrapeRequest:
			url = p.URL
		case *apiCrawlRequest:
			url = p.URL
		case *apiMapRequest:
			url = p.URL
		}
		if strings.TrimSpace(url) == "" {
			return apiFailure(ctx, 400, "url is required")
		}
		return h.call(ctx, operation, "", payload)
	}
}

func (h APISearchHandler) crawlStatus(ctx context.Context) error {
	id := strings.TrimSpace(anduin.Request(ctx).PathValue("id"))
	if id == "" {
		return apiFailure(ctx, 400, "crawl id is required")
	}
	return h.call(ctx, "crawl status", id, nil)
}

func (h APISearchHandler) call(ctx context.Context, operation, id string, payload any) error {
	result, err := h.client.Request(ctx, operation, id, payload)
	if err != nil {
		slog.ErrorContext(ctx, operation+" provider request failed", "error", err)
		return apiFailure(ctx, http.StatusBadGateway, operation+" request failed")
	}
	var response any
	switch operation {
	case "search":
		response = &apiSearchResponse{}
	case "scrape":
		response = &apiScrapeResponse{}
	case "crawl":
		response = &apiCrawlResponse{}
	case "crawl status":
		response = &apiCrawlStatusResponse{}
	case "map":
		response = &apiMapResponse{}
	}
	if err := json.Unmarshal(result, response); err != nil {
		slog.ErrorContext(ctx, operation+" provider response invalid", "error", err)
		return apiFailure(ctx, http.StatusBadGateway, operation+" request failed")
	}
	status := http.StatusOK
	if operation == "crawl" {
		status = http.StatusAccepted
	}
	return anduin.JSON(ctx, status, response)
}

func apiFailure(ctx context.Context, status int, message string) error {
	return anduin.JSON(ctx, status, apiError{Error: message})
}
