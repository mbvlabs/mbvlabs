package routes

import (
	"github.com/mbvlabs/andurel/pkg/routing"
)

const ProjectPrefix = "/projects"

var ProjectIndex = routing.NewSimpleRoute(
	"",
	"projects.index",
	ProjectPrefix,
	routing.TypeScriptRoute(),
)
var ProjectShow = routing.NewRouteWithSlug(
	"/:slug",
	"projects.show",
	ProjectPrefix,
	routing.TypeScriptRoute(),
)
