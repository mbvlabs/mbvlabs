package routes

import (
	"mbvlabs/internal/routing"
)

const WorkPrefix = "/work"

var WorkIndex = routing.NewSimpleRoute(
	"",
	"works.index",
	WorkPrefix,
)

var WorkShow = routing.NewRouteWithSlug(
	"/:slug",
	"works.show",
	WorkPrefix,
)
