package router

import (
	"net/http"

	"mbvlabs/controllers"
	"mbvlabs/router/routes"

	"github.com/labstack/echo/v4"
)

func registerWorkItemsRoutes(handler *echo.Echo, work_items controllers.WorkItems) {
	handler.Add(
		http.MethodGet, routes.WorkItemIndex.Path(), work_items.Index,
	).Name = routes.WorkItemIndex.Name()

	handler.Add(
		http.MethodGet, routes.WorkItemShow.Path(), work_items.Show,
	).Name = routes.WorkItemShow.Name()

	handler.Add(
		http.MethodGet, routes.WorkItemNew.Path(), work_items.New,
	).Name = routes.WorkItemNew.Name()

	handler.Add(
		http.MethodPost, routes.WorkItemCreate.Path(), work_items.Create,
	).Name = routes.WorkItemCreate.Name()

	handler.Add(
		http.MethodGet, routes.WorkItemEdit.Path(), work_items.Edit,
	).Name = routes.WorkItemEdit.Name()

	handler.Add(
		http.MethodPut, routes.WorkItemUpdate.Path(), work_items.Update,
	).Name = routes.WorkItemUpdate.Name()

	handler.Add(
		http.MethodDelete, routes.WorkItemDestroy.Path(), work_items.Destroy,
	).Name = routes.WorkItemDestroy.Name()
}
