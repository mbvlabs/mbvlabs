package controllers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"mbvlabs/internal/hypermedia"
	"mbvlabs/internal/storage"
	"mbvlabs/models"
	"mbvlabs/router"
	"mbvlabs/router/routes"
	"mbvlabs/views"

	"github.com/labstack/echo/v5"
)

type Works struct {
	db storage.Pool
}

func NewWorks(db storage.Pool) Works {
	return Works{db}
}

func (w Works) RegisterRoutes(r *router.Router) error {
	var errs []error
	var err error
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.WorkIndex.Path(),
		Name:    routes.WorkIndex.Name(),
		Handler: w.Index,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.WorkShow.Path(),
		Name:    routes.WorkShow.Name(),
		Handler: w.Show,
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (w Works) Index(etx *echo.Context) error {
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

	worksList, err := models.Work.PaginatePublished(
		etx.Request().Context(),
		w.db.Executor(),
		page,
		perPage,
	)
	if err != nil {
		slog.ErrorContext(etx.Request().Context(), "paginate published works", "error", err)
		return hypermedia.RenderPage(etx, views.InternalError())
	}

	return hypermedia.RenderPage(etx, views.WorkIndex{Items: worksList.Works}.Page())
}

func (w Works) Show(etx *echo.Context) error {
	slug := etx.Param("slug")
	if slug == "" {
		return hypermedia.RenderPage(etx, views.BadRequest())
	}

	work, err := models.Work.FindBySlug(etx.Request().Context(), w.db.Executor(), slug)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return hypermedia.RenderPage(etx, views.NotFound())
		}
		slog.ErrorContext(etx.Request().Context(), "find work by slug", "slug", slug, "error", err)
		return hypermedia.RenderPage(etx, views.InternalError())
	}
	if work.Status != models.Published {
		return hypermedia.RenderPage(etx, views.NotFound())
	}

	return hypermedia.RenderPage(etx, views.WorkShow{Item: work}.Page())
}
