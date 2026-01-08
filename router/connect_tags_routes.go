package router

import (
	"net/http"

	"mbvlabs/controllers"
	"mbvlabs/router/routes"

	"github.com/labstack/echo/v4"
)

func registerTagsRoutes(handler *echo.Echo, tags controllers.Tags) {
	handler.Add(
		http.MethodGet, routes.TagIndex.Path(), tags.Index,
	).Name = routes.TagIndex.Name()

	handler.Add(
		http.MethodGet, routes.TagShow.Path(), tags.Show,
	).Name = routes.TagShow.Name()

	handler.Add(
		http.MethodGet, routes.TagNew.Path(), tags.New,
	).Name = routes.TagNew.Name()

	handler.Add(
		http.MethodPost, routes.TagCreate.Path(), tags.Create,
	).Name = routes.TagCreate.Name()

	handler.Add(
		http.MethodGet, routes.TagEdit.Path(), tags.Edit,
	).Name = routes.TagEdit.Name()

	handler.Add(
		http.MethodPut, routes.TagUpdate.Path(), tags.Update,
	).Name = routes.TagUpdate.Name()

	handler.Add(
		http.MethodDelete, routes.TagDestroy.Path(), tags.Destroy,
	).Name = routes.TagDestroy.Name()
}
