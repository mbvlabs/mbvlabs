package router

import (
	"net/http"

	"mbvlabs/controllers"
	"mbvlabs/router/routes"

	"github.com/labstack/echo/v4"
)

func registerAPIRoutes(handler *echo.Echo, apiController controllers.API) {
	handler.Add(
		http.MethodGet, routes.Health.Path(), apiController.Health,
	).Name = routes.Health.Name()
}
