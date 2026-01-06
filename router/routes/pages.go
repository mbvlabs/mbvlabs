package routes

import (
	"mbvlabs/internal/routing"
)

const PagePrefix = "pages"

var HomePage = routing.NewSimpleRoute(
	"",
	"home",
	"",
)
