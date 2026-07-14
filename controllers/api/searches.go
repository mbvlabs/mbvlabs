package api

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"mbvlabs/config"
	"mbvlabs/router"
	"mbvlabs/router/middleware"
	"mbvlabs/router/routes"
	"mbvlabs/services"

	"github.com/labstack/echo/v5"
)

type Searches struct {
	search services.Search
	cfg    config.Config
}

func NewSearches(search services.Search, cfg config.Config) Searches {
	return Searches{search: search, cfg: cfg}
}

func (s Searches) RegisterRoutes(r *router.Router) error {
	auth := []echo.MiddlewareFunc{middleware.APIBasicAuth(s.cfg)}
	registered := []echo.Route{
		{
			Method:      http.MethodPost,
			Path:        routes.ApiSearch.Path(),
			Name:        routes.ApiSearch.Name(),
			Handler:     s.Search,
			Middlewares: auth,
		},
		{
			Method:      http.MethodPost,
			Path:        routes.ApiScrape.Path(),
			Name:        routes.ApiScrape.Name(),
			Handler:     s.Scrape,
			Middlewares: auth,
		},
		{
			Method:      http.MethodPost,
			Path:        routes.ApiCrawl.Path(),
			Name:        routes.ApiCrawl.Name(),
			Handler:     s.Crawl,
			Middlewares: auth,
		},
		{
			Method:      http.MethodGet,
			Path:        routes.ApiCrawlStatus.Path(),
			Name:        routes.ApiCrawlStatus.Name(),
			Handler:     s.CrawlStatus,
			Middlewares: auth,
		},
		{
			Method:      http.MethodPost,
			Path:        routes.ApiMap.Path(),
			Name:        routes.ApiMap.Name(),
			Handler:     s.Map,
			Middlewares: auth,
		},
	}

	var errs []error
	for _, route := range registered {
		if _, err := r.AddRoute(route); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (s Searches) Search(etx *echo.Context) error {
	var payload services.WebSearchRequest
	if err := etx.Bind(&payload); err != nil {
		return etx.JSON(http.StatusBadRequest, errorResponse{Error: "invalid JSON request body"})
	}
	if strings.TrimSpace(payload.Query) == "" {
		return etx.JSON(http.StatusBadRequest, errorResponse{Error: "query is required"})
	}

	response, err := s.search.Web(etx.Request().Context(), payload)
	if err != nil {
		return upstreamError(etx, "search", err)
	}
	return etx.JSON(http.StatusOK, response)
}

func (s Searches) Scrape(etx *echo.Context) error {
	var payload services.ScrapeRequest
	if err := etx.Bind(&payload); err != nil {
		return etx.JSON(http.StatusBadRequest, errorResponse{Error: "invalid JSON request body"})
	}
	if strings.TrimSpace(payload.URL) == "" {
		return etx.JSON(http.StatusBadRequest, errorResponse{Error: "url is required"})
	}

	response, err := s.search.Scrape(etx.Request().Context(), payload)
	if err != nil {
		return upstreamError(etx, "scrape", err)
	}
	return etx.JSON(http.StatusOK, response)
}

func (s Searches) Crawl(etx *echo.Context) error {
	var payload services.CrawlRequest
	if err := etx.Bind(&payload); err != nil {
		return etx.JSON(http.StatusBadRequest, errorResponse{Error: "invalid JSON request body"})
	}
	if strings.TrimSpace(payload.URL) == "" {
		return etx.JSON(http.StatusBadRequest, errorResponse{Error: "url is required"})
	}

	response, err := s.search.Crawl(etx.Request().Context(), payload)
	if err != nil {
		return upstreamError(etx, "crawl", err)
	}
	return etx.JSON(http.StatusAccepted, response)
}

func (s Searches) CrawlStatus(etx *echo.Context) error {
	id := strings.TrimSpace(etx.Param("id"))
	if id == "" {
		return etx.JSON(http.StatusBadRequest, errorResponse{Error: "crawl id is required"})
	}

	response, err := s.search.CrawlStatus(etx.Request().Context(), id)
	if err != nil {
		return upstreamError(etx, "crawl status", err)
	}
	return etx.JSON(http.StatusOK, response)
}

func (s Searches) Map(etx *echo.Context) error {
	var payload services.MapRequest
	if err := etx.Bind(&payload); err != nil {
		return etx.JSON(http.StatusBadRequest, errorResponse{Error: "invalid JSON request body"})
	}
	if strings.TrimSpace(payload.URL) == "" {
		return etx.JSON(http.StatusBadRequest, errorResponse{Error: "url is required"})
	}

	response, err := s.search.Map(etx.Request().Context(), payload)
	if err != nil {
		return upstreamError(etx, "map", err)
	}
	return etx.JSON(http.StatusOK, response)
}

func upstreamError(etx *echo.Context, operation string, err error) error {
	slog.ErrorContext(etx.Request().Context(), operation+" provider request failed", "error", err)
	return etx.JSON(http.StatusBadGateway, errorResponse{Error: operation + " request failed"})
}
