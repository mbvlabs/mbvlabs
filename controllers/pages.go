package controllers

import (
	"errors"
	"net/http"

	"mbvlabs/internal/hypermedia"
	"mbvlabs/internal/storage"
	"mbvlabs/models"
	"mbvlabs/queue"
	"mbvlabs/router"
	"mbvlabs/router/middleware"
	"mbvlabs/router/routes"
	"mbvlabs/views"

	"github.com/labstack/echo/v5"
)

type Pages struct {
	db         storage.Pool
	insertOnly queue.InsertOnly
}

func NewPages(
	db storage.Pool,
	insertOnly queue.InsertOnly,
) Pages {
	return Pages{db, insertOnly}
}

func (p Pages) RegisterRoutes(r *router.Router) error {
	errs := []error{}

	_, err := r.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.HomePage.Path(),
		Name:        routes.HomePage.Name(),
		Handler:     p.Home,
		Middlewares: []echo.MiddlewareFunc{middleware.MarkdownNegotiation},
	})
	if err != nil {
		errs = append(errs, err)
	}

	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodHead,
		Path:        routes.HomePage.Path(),
		Name:        routes.HomePage.Name() + ".head",
		Handler:     p.Home,
		Middlewares: []echo.MiddlewareFunc{middleware.MarkdownNegotiation},
	})
	if err != nil {
		errs = append(errs, err)
	}

	_ = r.AddRouteNotFound(p.NotFound)

	_, err = r.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        routes.AboutMe.Path(),
		Name:        routes.AboutMe.Name(),
		Handler:     p.AboutMe,
		Middlewares: []echo.MiddlewareFunc{middleware.MarkdownNegotiation},
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.LegacyAboutMe.Path(),
		Name:    routes.LegacyAboutMe.Name(),
		Handler: p.RedirectToAboutMe,
	})
	if err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func (p Pages) Home(etx *echo.Context) error {
	etx.Response().Header().Set("Link", "<"+routes.Sitemap.URL()+">; rel=\"describedby\"; type=\"application/xml\"")
	ctx := etx.Request().Context()

	featuredWorks, err := models.Work.FeaturedPublished(ctx, p.db.Executor(), 3)
	if err != nil {
		return hypermedia.RenderPage(etx, views.InternalError())
	}

	featuredProjects, err := models.Project.FeaturedPublished(ctx, p.db.Executor(), 3)
	if err != nil {
		return hypermedia.RenderPage(etx, views.InternalError())
	}

	return hypermedia.RenderPage(etx, views.LandingPage{
		FeaturedWorks:    featuredWorks,
		FeaturedProjects: featuredProjects,
	}.Page())
}

func (p Pages) NotFound(etx *echo.Context) error {
	return hypermedia.RenderPage(etx, views.NotFound())
}

func (p Pages) AboutMe(etx *echo.Context) error {
	return hypermedia.RenderPage(etx, views.AboutMe())
}

func (p Pages) RedirectToAboutMe(etx *echo.Context) error {
	return etx.Redirect(http.StatusMovedPermanently, routes.AboutMe.URL())
}
