package controllers

import (
	"errors"
	"mbvlabs/internal/hypermedia"
	"mbvlabs/router"
	"mbvlabs/router/routes"
	"mbvlabs/views"
	"net/http"

	"github.com/labstack/echo/v5"
)

type ServiceOfferings struct{}

func NewServiceOfferings() ServiceOfferings {
	return ServiceOfferings{}
}

func (so ServiceOfferings) Index(etx *echo.Context) error {
	return hypermedia.RenderPage(
		etx,
		views.ServiceOfferingIndex{Items: []views.ServiceOfferingData{}}.Page(),
	)
}

func (so ServiceOfferings) RegisterRoutes(r *router.Router) error {
	var errs []error
	var err error
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.ServiceOfferingIndex.Path(),
		Name:    routes.ServiceOfferingIndex.Name(),
		Handler: so.Index,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.LegacyServiceOfferingIndex.Path(),
		Name:    routes.LegacyServiceOfferingIndex.Name(),
		Handler: so.RedirectToIndex,
	})
	if err != nil {
		errs = append(errs, err)
	}
	_, err = r.AddRoute(echo.Route{
		Method:  http.MethodGet,
		Path:    routes.LegacyServiceOfferingShow.Path(),
		Name:    routes.LegacyServiceOfferingShow.Name(),
		Handler: so.RedirectToIndex,
	})
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (so ServiceOfferings) RedirectToIndex(etx *echo.Context) error {
	return etx.Redirect(http.StatusMovedPermanently, routes.ServiceOfferingIndex.URL())
}
