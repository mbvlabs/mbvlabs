// Package siteapi defines the HTTP integration payloads whose open-ended JSON
// objects are not consumed by the Inertia frontend.
package siteapi

import "encoding/json"

type APIError struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

type WebSearchRequest struct {
	Query        string `json:"query"`
	Location     string `json:"location,omitempty"`
	CountryCode  string `json:"countryCode,omitempty"`
	LanguageCode string `json:"languageCode,omitempty"`
	Autocorrect  *bool  `json:"autocorrect,omitempty"`
	Page         int    `json:"page,omitempty"`
	Num          int    `json:"num,omitempty"`
	TimeRange    string `json:"timeRange,omitempty"`
}

// APISearchProviderRequest is the translated request sent to the search provider.
type APISearchProviderRequest struct {
	Query        string `json:"q"`
	Location     string `json:"location,omitempty"`
	CountryCode  string `json:"gl,omitempty"`
	LanguageCode string `json:"hl,omitempty"`
	Autocorrect  *bool  `json:"autocorrect,omitempty"`
	Page         int    `json:"page,omitempty"`
	Num          int    `json:"num,omitempty"`
	TimeRange    string `json:"tbs,omitempty"`
}

type APILocation struct {
	Country   string   `json:"country,omitempty"`
	Languages []string `json:"languages,omitempty"`
}
type APIParser struct {
	Type     string `json:"type"`
	Mode     string `json:"mode,omitempty"`
	MaxPages int    `json:"maxPages,omitempty"`
}
type APIProfile struct {
	Name        string `json:"name"`
	SaveChanges *bool  `json:"saveChanges,omitempty"`
}
type APIScrapeOptions struct {
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
	Parsers             []APIParser       `json:"parsers,omitempty"`
	Actions             []map[string]any  `json:"actions,omitempty"`
	Location            *APILocation      `json:"location,omitempty"`
	RemoveBase64Images  *bool             `json:"removeBase64Images,omitempty"`
	BlockAds            *bool             `json:"blockAds,omitempty"`
	Proxy               string            `json:"proxy,omitempty"`
	StoreInCache        *bool             `json:"storeInCache,omitempty"`
	Lockdown            bool              `json:"lockdown,omitempty"`
	RedactPII           any               `json:"redactPII,omitempty"`
	Profile             *APIProfile       `json:"profile,omitempty"`
	ThreatProtection    map[string]any    `json:"threatProtection,omitempty"`
}
type APIScrapeRequest struct {
	URL string `json:"url"`
	APIScrapeOptions
	ZeroDataRetention bool `json:"zeroDataRetention,omitempty"`
}
type APIWebhook struct {
	URL      string            `json:"url"`
	Headers  map[string]string `json:"headers,omitempty"`
	Metadata map[string]any    `json:"metadata,omitempty"`
	Events   []string          `json:"events,omitempty"`
}
type APICrawlRequest struct {
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
	Webhook               *APIWebhook       `json:"webhook,omitempty"`
	ScrapeOptions         *APIScrapeOptions `json:"scrapeOptions,omitempty"`
	ZeroDataRetention     bool              `json:"zeroDataRetention,omitempty"`
}
type APIMapRequest struct {
	URL                   string         `json:"url"`
	Search                string         `json:"search,omitempty"`
	Sitemap               string         `json:"sitemap,omitempty"`
	IncludeSubdomains     *bool          `json:"includeSubdomains,omitempty"`
	IgnoreQueryParameters *bool          `json:"ignoreQueryParameters,omitempty"`
	IgnoreCache           bool           `json:"ignoreCache,omitempty"`
	Limit                 int            `json:"limit,omitempty"`
	Timeout               int            `json:"timeout,omitempty"`
	Location              *APILocation   `json:"location,omitempty"`
	ThreatProtection      map[string]any `json:"threatProtection,omitempty"`
}

type APIDocument struct {
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
type APISearchResponse struct {
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
type APIScrapeResponse struct {
	Success bool        `json:"success"`
	Data    APIDocument `json:"data"`
}
type APICrawlResponse struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
	URL     string `json:"url"`
}
type APICrawlStatusResponse struct {
	Status      string        `json:"status"`
	Total       int           `json:"total"`
	Completed   int           `json:"completed"`
	CreditsUsed int           `json:"creditsUsed"`
	ExpiresAt   string        `json:"expiresAt"`
	CreatedAt   string        `json:"createdAt"`
	CompletedAt string        `json:"completedAt,omitempty"`
	Duration    float64       `json:"duration"`
	Next        string        `json:"next,omitempty"`
	Data        []APIDocument `json:"data"`
}
type APIMapResponse struct {
	Success bool `json:"success"`
	Links   []struct {
		URL         string `json:"url"`
		Title       string `json:"title,omitempty"`
		Description string `json:"description,omitempty"`
	} `json:"links"`
}
