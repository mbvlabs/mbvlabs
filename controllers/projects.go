package controllers

import (
	"errors"
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

type Projects struct {
	db storage.Pool
}

func NewProjects(db storage.Pool) Projects {
	return Projects{db}
}

func (p Projects) Index(etx *echo.Context) error {
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

	projectsList, err := models.Project.PaginatePublished(
		etx.Request().Context(),
		p.db.Executor(),
		page,
		perPage,
	)
	if err != nil {
		return hypermedia.RenderPage(etx, views.InternalError())
	}

	return hypermedia.RenderPage(etx, views.ProjectIndex{Items: projectsList.Projects}.Page())
}

func (p Projects) Show(etx *echo.Context) error {
	slug := etx.Param("slug")
	if slug == "" {
		return hypermedia.RenderPage(etx, views.BadRequest())
	}

	project, err := models.Project.FindBySlug(etx.Request().Context(), p.db.Executor(), slug)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return hypermedia.RenderPage(etx, views.NotFound())
		}
		return hypermedia.RenderPage(etx, views.InternalError())
	}
	if !project.PublishedAt.Valid {
		return hypermedia.RenderPage(etx, views.NotFound())
	}

	return hypermedia.RenderPage(etx, views.ProjectShow{Item: project}.Page())
}

func (p Projects) RegisterRoutes(r *router.Router) error {
	var errs []error
	var err error
	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.ProjectIndex.Path(),
		Name:        routes.ProjectIndex.Name(),
		Handler:     p.Index,
		Middlewares: []echo.MiddlewareFunc{middleware.MarkdownNegotiation},
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.ProjectShow.Path(),
		Name:        routes.ProjectShow.Name(),
		Handler:     p.Show,
		Middlewares: []echo.MiddlewareFunc{middleware.MarkdownNegotiation},
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
