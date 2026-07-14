package services

import (
	"context"

	"mbvlabs/clients/firecrawl"
	"mbvlabs/clients/serper"
)

type Search struct {
	serper    *serper.Client
	firecrawl *firecrawl.Client
}

func NewSearch(serperClient *serper.Client, firecrawlClient *firecrawl.Client) Search {
	return Search{serper: serperClient, firecrawl: firecrawlClient}
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

type WebSearchResponse = serper.SearchResponse
type ScrapeRequest = firecrawl.ScrapeRequest
type ScrapeResponse = firecrawl.ScrapeResponse
type CrawlRequest = firecrawl.CrawlRequest
type CrawlResponse = firecrawl.CrawlResponse
type CrawlStatusResponse = firecrawl.CrawlStatusResponse
type MapRequest = firecrawl.MapRequest
type MapResponse = firecrawl.MapResponse

func (s Search) Web(ctx context.Context, request WebSearchRequest) (WebSearchResponse, error) {
	return s.serper.Search(ctx, serper.Query{
		Query:        request.Query,
		Location:     request.Location,
		CountryCode:  request.CountryCode,
		LanguageCode: request.LanguageCode,
		Autocorrect:  request.Autocorrect,
		Page:         request.Page,
		Num:          request.Num,
		TimeRange:    request.TimeRange,
	})
}

func (s Search) Scrape(ctx context.Context, request ScrapeRequest) (ScrapeResponse, error) {
	return s.firecrawl.Scrape(ctx, request)
}

func (s Search) Crawl(ctx context.Context, request CrawlRequest) (CrawlResponse, error) {
	return s.firecrawl.Crawl(ctx, request)
}

func (s Search) CrawlStatus(ctx context.Context, id string) (CrawlStatusResponse, error) {
	return s.firecrawl.CrawlStatus(ctx, id)
}

func (s Search) Map(ctx context.Context, request MapRequest) (MapResponse, error) {
	return s.firecrawl.Map(ctx, request)
}
