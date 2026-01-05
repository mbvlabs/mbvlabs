package routes

import (
	"mbvlabs/internal/routing"
)

const APIPrefix = "api"

var Health = routing.NewSimpleRoute(
	"/health",
	"health",
	APIPrefix,
)
