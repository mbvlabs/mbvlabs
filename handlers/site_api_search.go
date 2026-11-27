package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"mbvlabs/middleware"
	"mbvlabs/payloads/siteapi"
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

func (h APISearchHandler) search(ctx context.Context) error {
	var p siteapi.WebSearchRequest
	if err := anduin.DecodeJSON(ctx, &p); err != nil {
		return apiFailure(ctx, 400, "invalid JSON request body")
	}
	if strings.TrimSpace(p.Query) == "" {
		return apiFailure(ctx, 400, "query is required")
	}
	request := siteapi.APISearchProviderRequest{Query: p.Query, Location: p.Location, CountryCode: p.CountryCode, LanguageCode: p.LanguageCode, Autocorrect: p.Autocorrect, Page: p.Page, Num: p.Num, TimeRange: p.TimeRange}
	return h.call(ctx, "search", "", request)
}

func (h APISearchHandler) post(operation string) anduin.Handler {
	return func(ctx context.Context) error {
		var payload any
		switch operation {
		case "scrape":
			payload = &siteapi.APIScrapeRequest{}
		case "crawl":
			payload = &siteapi.APICrawlRequest{}
		case "map":
			payload = &siteapi.APIMapRequest{}
		}
		if err := anduin.DecodeJSON(ctx, payload); err != nil {
			return apiFailure(ctx, 400, "invalid JSON request body")
		}
		var url string
		switch p := payload.(type) {
		case *siteapi.APIScrapeRequest:
			url = p.URL
		case *siteapi.APICrawlRequest:
			url = p.URL
		case *siteapi.APIMapRequest:
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
		response = &siteapi.APISearchResponse{}
	case "scrape":
		response = &siteapi.APIScrapeResponse{}
	case "crawl":
		response = &siteapi.APICrawlResponse{}
	case "crawl status":
		response = &siteapi.APICrawlStatusResponse{}
	case "map":
		response = &siteapi.APIMapResponse{}
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
	return anduin.JSON(ctx, status, siteapi.APIError{Error: message})
}
