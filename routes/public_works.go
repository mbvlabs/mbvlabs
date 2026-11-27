package routes

import (
	"github.com/mbvlabs/andurel/pkg/routing"
)

const WorkPrefix = "/work"

var WorkIndex = routing.NewSimpleRoute(
	"",
	"works.index",
	WorkPrefix,
	routing.TypeScriptRoute(),
)

var WorkShow = routing.NewRouteWithSlug(
	"/:slug",
	"works.show",
	WorkPrefix,
	routing.TypeScriptRoute(),
)
