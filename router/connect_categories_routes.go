package router

import (
	"net/http"

	"mbvlabs/controllers"
	"mbvlabs/router/routes"

	"github.com/labstack/echo/v4"
)

func registerCategoriesRoutes(handler *echo.Echo, categories controllers.Categorys) {
	handler.Add(
		http.MethodGet, routes.CategoryIndex.Path(), categories.Index,
	).Name = routes.CategoryIndex.Name()

	handler.Add(
		http.MethodGet, routes.CategoryShow.Path(), categories.Show,
	).Name = routes.CategoryShow.Name()

	handler.Add(
		http.MethodGet, routes.CategoryNew.Path(), categories.New,
	).Name = routes.CategoryNew.Name()

	handler.Add(
		http.MethodPost, routes.CategoryCreate.Path(), categories.Create,
	).Name = routes.CategoryCreate.Name()

	handler.Add(
		http.MethodGet, routes.CategoryEdit.Path(), categories.Edit,
	).Name = routes.CategoryEdit.Name()

	handler.Add(
		http.MethodPut, routes.CategoryUpdate.Path(), categories.Update,
	).Name = routes.CategoryUpdate.Name()

	handler.Add(
		http.MethodDelete, routes.CategoryDestroy.Path(), categories.Destroy,
	).Name = routes.CategoryDestroy.Name()
}
