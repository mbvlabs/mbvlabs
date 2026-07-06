package routes

import (
	"mbvlabs/internal/routing"
)

const APIPrefix = "/api"

var Health = routing.NewSimpleRoute(
	"/health",
	"api.health",
	APIPrefix,
)
