package routes

import (
	"mbvlabs/internal/routing"
)

const CategoryPrefix = "categories"

var CategoryIndex = routing.NewSimpleRoute(
	"/",
	"index",
	CategoryPrefix,
)

var CategoryShow = routing.NewRouteWithID(
	"/:id",
	"show",
	CategoryPrefix,
)

var CategoryNew = routing.NewSimpleRoute(
	"/new",
	"new",
	CategoryPrefix,
)

var CategoryCreate = routing.NewSimpleRoute(
	"/",
	"create",
	CategoryPrefix,
)

var CategoryEdit = routing.NewRouteWithID(
	"/:id/edit",
	"edit",
	CategoryPrefix,
)

var CategoryUpdate = routing.NewRouteWithID(
	"/:id",
	"update",
	CategoryPrefix,
)

var CategoryDestroy = routing.NewRouteWithID(
	"/:id",
	"destroy",
	CategoryPrefix,
)
