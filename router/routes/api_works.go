package routes

import (
	"mbvlabs/internal/routing"
)

const ApiWorkPrefix = "/api/works"

var ApiWorkCreate = routing.NewSimpleRoute(
	"",
	"api.works.create",
	ApiWorkPrefix,
)
