package routes

import "mbvlabs/internal/routing"

const AdminPrefix = "/admin"

var AdminHomePage = routing.NewSimpleRoute(
	"",
	"admin.home",
	AdminPrefix,
)
