package routes

import (
	"mbvlabs/internal/routing"
)

const WorkItemPrefix = "work_items"

var WorkItemIndex = routing.NewSimpleRoute(
	"/",
	"index",
	WorkItemPrefix,
)

var WorkItemShow = routing.NewRouteWithID(
	"/:id",
	"show",
	WorkItemPrefix,
)

var WorkItemNew = routing.NewSimpleRoute(
	"/new",
	"new",
	WorkItemPrefix,
)

var WorkItemCreate = routing.NewSimpleRoute(
	"/",
	"create",
	WorkItemPrefix,
)

var WorkItemEdit = routing.NewRouteWithID(
	"/:id/edit",
	"edit",
	WorkItemPrefix,
)

var WorkItemUpdate = routing.NewRouteWithID(
	"/:id",
	"update",
	WorkItemPrefix,
)

var WorkItemDestroy = routing.NewRouteWithID(
	"/:id",
	"destroy",
	WorkItemPrefix,
)
