package api

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"mbvlabs/config"
	"mbvlabs/internal/storage"
	"mbvlabs/internal/validation"
	"mbvlabs/models"
	"mbvlabs/router"
	"mbvlabs/router/middleware"
	"mbvlabs/router/routes"

	"github.com/gosimple/slug"
	"github.com/labstack/echo/v5"
)

type Articles struct {
	db  storage.Pool
	cfg config.Config
}

func NewArticles(db storage.Pool, cfg config.Config) Articles {
	return Articles{db: db, cfg: cfg}
}

func (a Articles) RegisterRoutes(r *router.Router) error {
	_, err := r.AddRoute(echo.Route{
		Method:  http.MethodPost,
		Path:    routes.ApiArticleCreate.Path(),
		Name:    routes.ApiArticleCreate.Name(),
		Handler: a.Create,
		Middlewares: []echo.MiddlewareFunc{
			middleware.APIBasicAuth(a.cfg),
		},
	})
	return err
}

type CreateArticlePayload struct {
	Title         string   `json:"title"`
	Slug          string   `json:"slug,omitempty"`
	Excerpt       string   `json:"excerpt,omitempty"`
	Body          string   `json:"body,omitempty"`
	CoverImageURL string   `json:"coverImageUrl,omitempty"`
	Tags          []string `json:"tags,omitempty"`
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
	article, err := models.BlogPost.Create(
		etx.Request().Context(),
		a.db.Executor(),
		models.CreateBlogPostData{
			Title:         payload.Title,
			Slug:          slug.Make(slugSource),
			Excerpt:       payload.Excerpt,
			Body:          payload.Body,
			Status:        models.Draft.String(),
			CoverImageUrl: payload.CoverImageURL,
			Tags:          payload.Tags,
		},
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
