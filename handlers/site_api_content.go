package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"mbvlabs/domains"
	"mbvlabs/middleware"
	"mbvlabs/routes"

	"github.com/gosimple/slug"
	"github.com/mbvlabs/andurel/pkg/anduin"
	"github.com/mbvlabs/andurel/pkg/validation"
)

type apiArticleStore interface {
	AllDraftOrPublished(context.Context) ([]domains.BlogPostEntity, error)
}
type apiArticleCreator interface {
	Create(context.Context, domains.CreateBlogPostData, *time.Time) (domains.BlogPostEntity, error)
}
type apiWorkCreator interface {
	CreateWork(context.Context, domains.CreateWorkData) (domains.WorkEntity, error)
}

type APIContentHandler struct {
	articles           apiArticleStore
	articleCreator     apiArticleCreator
	works              apiWorkCreator
	featuredLimit      error
	username, password string
}

// featuredLimit is the work service's sentinel error for its three-item limit.
func NewAPIContentHandler(articles apiArticleStore, creator apiArticleCreator, works apiWorkCreator, featuredLimit error, username, password string) APIContentHandler {
	return APIContentHandler{articles, creator, works, featuredLimit, username, password}
}

func (h APIContentHandler) Register(rt *anduin.Router) {
	auth := middleware.APIBasicAuth(h.username, h.password)
	rt.Route(routes.ApiArticleIndex).Get(h.articleIndex, auth)
	rt.Route(routes.ApiArticleCreate).Post(h.articleCreate, auth)
	// The legacy work creation endpoint is intentionally public.
	rt.Route(routes.ApiWorkCreate).Post(h.workCreate)
}

type createArticlePayload struct {
	Title         string   `json:"title"`
	Slug          string   `json:"slug"`
	Excerpt       string   `json:"excerpt"`
	Body          string   `json:"body"`
	CoverImageURL string   `json:"coverImageUrl"`
	Tags          []string `json:"tags"`
	ScheduledAt   *string  `json:"scheduledAt"`
}

func (h APIContentHandler) articleIndex(ctx context.Context) error {
	articles, err := h.articles.AllDraftOrPublished(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to list articles", "error", err)
		return apiFailure(ctx, 500, "failed to list articles")
	}
	return anduin.JSON(ctx, 200, articles)
}

func (h APIContentHandler) articleCreate(ctx context.Context) error {
	var p createArticlePayload
	if err := anduin.DecodeJSON(ctx, &p); err != nil {
		return apiFailure(ctx, 400, "invalid JSON request body")
	}
	source := strings.TrimSpace(p.Slug)
	if source == "" {
		source = p.Title
	}
	var scheduled *time.Time
	if p.ScheduledAt != nil {
		value, err := time.Parse(time.RFC3339, *p.ScheduledAt)
		if err != nil {
			return anduin.JSON(ctx, 400, apiError{"validation failed", map[string]string{"scheduledAt": "must be a valid RFC3339 timestamp"}})
		}
		scheduled = &value
	}
	post, err := h.articleCreator.Create(ctx, domains.CreateBlogPostData{
		Title: p.Title, Slug: slug.Make(source), Excerpt: p.Excerpt, Body: p.Body,
		Status: domains.Draft.String(), CoverImageUrl: p.CoverImageURL, Tags: p.Tags,
	}, scheduled)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create draft article", "error", err)
		if fields, ok := validation.As(err); ok {
			return anduin.JSON(ctx, 400, apiError{"validation failed", fields.ToMap()})
		}
		if domains.IsDuplicateBlogPostSlug(err) {
			return anduin.JSON(ctx, 409, apiError{"article slug already exists", map[string]string{"slug": "must be unique"}})
		}
		if errors.Is(err, domains.ErrDomainValidation) {
			return apiFailure(ctx, 400, "validation failed")
		}
		return apiFailure(ctx, 500, "failed to create article")
	}
	return anduin.JSON(ctx, http.StatusCreated, post)
}

type createWorkPayload struct {
	Title          string   `json:"title"`
	Slug           string   `json:"slug"`
	ClientName     string   `json:"clientName"`
	ClientIndustry string   `json:"clientIndustry"`
	ClientURL      string   `json:"clientUrl"`
	ClientLogoURL  string   `json:"clientLogoUrl"`
	Summary        string   `json:"summary"`
	CoverImageURL  string   `json:"coverImageUrl"`
	Specialisms    []string `json:"specialisms"`
	Platforms      []string `json:"platforms"`
	Technologies   []string `json:"technologies"`
	Challenge      string   `json:"challenge"`
	Approach       string   `json:"approach"`
	Deliverables   string   `json:"deliverables"`
	Outcome        string   `json:"outcome"`
	Content        string   `json:"content"`
	StartedAt      string   `json:"startedAt"`
	CompletedAt    string   `json:"completedAt"`
	Status         string   `json:"status"`
	PublishedAt    string   `json:"publishedAt"`
	IsFeatured     bool     `json:"isFeatured"`
}

func (h APIContentHandler) workCreate(ctx context.Context) error {
	var p createWorkPayload
	if err := anduin.DecodeJSON(ctx, &p); err != nil {
		return apiFailure(ctx, 400, "invalid JSON request body")
	}
	fields := map[string]string{}
	parse := func(field, value string) time.Time {
		if strings.TrimSpace(value) == "" {
			return time.Time{}
		}
		parsed, err := time.Parse("2006-01-02", value)
		if err != nil {
			fields[field] = "must be a valid date in YYYY-MM-DD format"
		}
		return parsed
	}
	start, end, published := parse("startedAt", p.StartedAt), parse("completedAt", p.CompletedAt), parse("publishedAt", p.PublishedAt)
	if len(fields) != 0 {
		return anduin.JSON(ctx, 400, apiError{"validation failed", fields})
	}
	status := strings.TrimSpace(p.Status)
	if status == "" {
		status = domains.Draft.String()
	}
	work, err := h.works.CreateWork(ctx, domains.CreateWorkData{
		Title: p.Title, Slug: p.Slug, ClientName: p.ClientName, ClientIndustry: p.ClientIndustry,
		ClientURL: p.ClientURL, ClientLogoURL: p.ClientLogoURL, Summary: p.Summary, CoverImageURL: p.CoverImageURL,
		Specialisms: p.Specialisms, Platforms: p.Platforms, Technologies: p.Technologies, Challenge: p.Challenge,
		Approach: p.Approach, Deliverables: p.Deliverables, Outcome: p.Outcome, Content: p.Content,
		StartedAt: start, CompletedAt: end, PublishedAt: published, Status: domains.StatusEnum(status), IsFeatured: p.IsFeatured,
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to create work", "error", err)
		if fields, ok := validation.As(err); ok {
			return anduin.JSON(ctx, 400, apiError{"validation failed", fields.ToMap()})
		}
		if h.featuredLimit != nil && errors.Is(err, h.featuredLimit) {
			return anduin.JSON(ctx, 400, apiError{"validation failed", map[string]string{"isFeatured": "No more than 3 works can be featured."}})
		}
		if domains.IsDuplicateWorkSlug(err) {
			return anduin.JSON(ctx, 409, apiError{"work slug already exists", map[string]string{"slug": "must be unique"}})
		}
		return apiFailure(ctx, 500, fmt.Sprintf("failed to create work: %v", err))
	}
	return anduin.JSON(ctx, http.StatusCreated, work)
}
