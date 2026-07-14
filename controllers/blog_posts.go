package controllers

import (
	"errors"
	"log/slog"
	"mbvlabs/internal/hypermedia"
	"mbvlabs/internal/storage"
	"mbvlabs/models"
	"mbvlabs/router"
	"mbvlabs/router/middleware"
	"mbvlabs/router/routes"
	"mbvlabs/views"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

type BlogPosts struct {
	db storage.Pool
}

func NewBlogPosts(db storage.Pool) BlogPosts {
	return BlogPosts{db}
}

func (bp BlogPosts) Index(etx *echo.Context) error {
	page := int64(1)
	if p := etx.QueryParam("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = int64(parsed)
		}
	}

	perPage := int64(25)
	if pp := etx.QueryParam("per_page"); pp != "" {
		if parsed, err := strconv.Atoi(pp); err == nil && parsed > 0 &&
			parsed <= 100 {
			perPage = int64(parsed)
		}
	}

	blogPostsList, err := models.BlogPost.PaginatePublished(
		etx.Request().Context(),
		bp.db.Executor(),
		page,
		perPage,
	)
	if err != nil {
		return hypermedia.RenderPage(etx, views.InternalError())
	}

	return hypermedia.RenderPage(etx, views.BlogPostIndex{Items: blogPostsList.BlogPosts}.Page())
}

func (bp BlogPosts) Show(etx *echo.Context) error {
	slug := etx.Param("slug")
	if slug == "" {
		return hypermedia.RenderPage(etx, views.BadRequest())
	}

	blogPost, err := models.BlogPost.FindBySlug(etx.Request().Context(), bp.db.Executor(), slug)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return hypermedia.RenderPage(etx, views.NotFound())
		}
		slog.ErrorContext(
			etx.Request().Context(),
			"find blog post by slug",
			"slug",
			slug,
			"error",
			err,
		)
		return hypermedia.RenderPage(etx, views.InternalError())
	}
	if blogPost.Status != "published" {
		return hypermedia.RenderPage(etx, views.NotFound())
	}

	return hypermedia.RenderPage(etx, views.BlogPostShow{Item: blogPost}.Page())
}

func (bp BlogPosts) RegisterRoutes(r *router.Router) error {
	var errs []error
	var err error
	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.BlogPostIndex.Path(),
		Name:        routes.BlogPostIndex.Name(),
		Handler:     bp.Index,
		Middlewares: []echo.MiddlewareFunc{middleware.MarkdownNegotiation},
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.BlogPostShow.Path(),
		Name:        routes.BlogPostShow.Name(),
		Handler:     bp.Show,
		Middlewares: []echo.MiddlewareFunc{middleware.MarkdownNegotiation},
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.LegacyBlogPostIndex.Path(),
		Name:    routes.LegacyBlogPostIndex.Name(),
		Handler: bp.RedirectToIndex,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.LegacyBlogPostShow.Path(),
		Name:    routes.LegacyBlogPostShow.Name(),
		Handler: bp.RedirectToShow,
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (bp BlogPosts) RedirectToIndex(etx *echo.Context) error {
	return etx.Redirect(http.StatusMovedPermanently, routes.BlogPostIndex.URL())
}

func (bp BlogPosts) RedirectToShow(etx *echo.Context) error {
	slug := etx.Param("slug")
	if slug == "" {
		return etx.Redirect(http.StatusMovedPermanently, routes.BlogPostIndex.URL())
	}
	return etx.Redirect(http.StatusMovedPermanently, routes.BlogPostShow.URL(slug))
}
