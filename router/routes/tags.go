package routes

import (
	"mbvlabs/internal/routing"
)

const TagPrefix = "tags"

var TagIndex = routing.NewSimpleRoute(
	"/",
	"index",
	TagPrefix,
)

var TagShow = routing.NewRouteWithID(
	"/:id",
	"show",
	TagPrefix,
)

var TagNew = routing.NewSimpleRoute(
	"/new",
	"new",
	TagPrefix,
)

var TagCreate = routing.NewSimpleRoute(
	"/",
	"create",
	TagPrefix,
)

var TagEdit = routing.NewRouteWithID(
	"/:id/edit",
	"edit",
	TagPrefix,
)

var TagUpdate = routing.NewRouteWithID(
	"/:id",
	"update",
	TagPrefix,
)

var TagDestroy = routing.NewRouteWithID(
	"/:id",
	"destroy",
	TagPrefix,
)
