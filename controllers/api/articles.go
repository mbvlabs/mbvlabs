package api

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"mbvlabs/config"
	"mbvlabs/internal/storage"
	"mbvlabs/internal/validation"
	"mbvlabs/models"
	"mbvlabs/router"
	"mbvlabs/router/middleware"
	"mbvlabs/router/routes"
	"mbvlabs/services"

	"github.com/gosimple/slug"
	"github.com/labstack/echo/v5"
)

type Articles struct {
	db  storage.Pool
	svc services.BlogPosts
	cfg config.Config
}

func NewArticles(db storage.Pool, svc services.BlogPosts, cfg config.Config) Articles {
	return Articles{db: db, svc: svc, cfg: cfg}
}

func (a Articles) RegisterRoutes(r *router.Router) error {
	errs := []error{}
	for _, route := range []echo.Route{
		{
			Method:  http.MethodGet,
			Path:    routes.ApiArticleIndex.Path(),
			Name:    routes.ApiArticleIndex.Name(),
			Handler: a.Index,
			Middlewares: []echo.MiddlewareFunc{
				middleware.APIBasicAuth(a.cfg),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    routes.ApiArticleCreate.Path(),
			Name:    routes.ApiArticleCreate.Name(),
			Handler: a.Create,
			Middlewares: []echo.MiddlewareFunc{
				middleware.APIBasicAuth(a.cfg),
			},
		},
	} {
		if _, err := r.AddRoute(route); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (a Articles) Index(etx *echo.Context) error {
	articles, err := models.BlogPost.AllDraftOrPublished(etx.Request().Context(), a.db.Executor())
	if err != nil {
		slog.ErrorContext(etx.Request().Context(), "failed to list articles", "error", err)
		return etx.JSON(http.StatusInternalServerError, errorResponse{Error: "failed to list articles"})
	}

	return etx.JSON(http.StatusOK, articles)
}

type CreateArticlePayload struct {
	Title         string   `json:"title"`
	Slug          string   `json:"slug,omitempty"`
	Excerpt       string   `json:"excerpt,omitempty"`
	Body          string   `json:"body,omitempty"`
	CoverImageURL string   `json:"coverImageUrl,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	ScheduledAt   *string  `json:"scheduledAt,omitempty"`
}

func (a Articles) Create(etx *echo.Context) error {
	var payload CreateArticlePayload
	if err := etx.Bind(&payload); err != nil {
		return etx.JSON(http.StatusBadRequest, errorResponse{Error: "invalid JSON request body"})
	}

	slugSource := strings.TrimSpace(payload.Slug)
	if slugSource == "" {
		slugSource = payload.Title
	}
	var scheduledAt *time.Time
	if payload.ScheduledAt != nil {
		parsed, err := time.Parse(time.RFC3339, *payload.ScheduledAt)
		if err != nil {
			return etx.JSON(http.StatusBadRequest, errorResponse{
				Error:  "validation failed",
				Fields: map[string]string{"scheduledAt": "must be a valid RFC3339 timestamp"},
			})
		}
		scheduledAt = &parsed
	}

	article, err := a.svc.Create(
		etx.Request().Context(),
		models.CreateBlogPostData{
			Title:         payload.Title,
			Slug:          slug.Make(slugSource),
			Excerpt:       payload.Excerpt,
			Body:          payload.Body,
			Status:        models.Draft.String(),
			CoverImageUrl: payload.CoverImageURL,
			Tags:          payload.Tags,
		},
		scheduledAt,
	)
	if err != nil {
		slog.ErrorContext(etx.Request().Context(), "failed to create draft article", "error", err)
		if fields, ok := validation.As(err); ok {
			return etx.JSON(http.StatusBadRequest, errorResponse{
				Error:  "validation failed",
				Fields: fields.ToMap(),
			})
		}
		if isUniqueViolation(err, "blog_posts_slug_key") {
			return etx.JSON(http.StatusConflict, errorResponse{
				Error:  "article slug already exists",
				Fields: map[string]string{"slug": "must be unique"},
			})
		}
		if errors.Is(err, models.ErrDomainValidation) {
			return etx.JSON(http.StatusBadRequest, errorResponse{Error: "validation failed"})
		}
		return etx.JSON(http.StatusInternalServerError, errorResponse{Error: "failed to create article"})
	}

	return etx.JSON(http.StatusCreated, article)
}
