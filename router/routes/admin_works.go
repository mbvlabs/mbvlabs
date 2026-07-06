package routes

import (
	"mbvlabs/internal/routing"
)

const AdminWorkPrefix = "/admin/works"

var AdminWorkIndex = routing.NewSimpleRoute(
	"",
	"admin.works.index",
	AdminWorkPrefix,
)
var AdminWorkShow = routing.NewRouteWithSerialID(
	"/:id",
	"admin.works.show",
	AdminWorkPrefix,
)
var AdminWorkNew = routing.NewSimpleRoute(
	"/new",
	"admin.works.new",
	AdminWorkPrefix,
)
var AdminWorkCreate = routing.NewSimpleRoute(
	"",
	"admin.works.create",
	AdminWorkPrefix,
)
var AdminWorkEdit = routing.NewRouteWithSerialID(
	"/:id/edit",
	"admin.works.edit",
	AdminWorkPrefix,
)
var AdminWorkUpdate = routing.NewRouteWithSerialID(
	"/:id",
	"admin.works.update",
	AdminWorkPrefix,
)
var AdminWorkDestroy = routing.NewRouteWithSerialID(
	"/:id",
	"admin.works.destroy",
	AdminWorkPrefix,
)
