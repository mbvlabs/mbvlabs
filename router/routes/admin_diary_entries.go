package routes

import (
	"mbvlabs/internal/routing"
)

const AdminDiaryEntryPrefix = "/admin/diary"

var AdminDiaryEntryIndex = routing.NewSimpleRoute(
	"",
	"admin.diary_entries.index",
	AdminDiaryEntryPrefix,
)
var AdminDiaryEntryShow = routing.NewRouteWithSerialID(
	"/:id",
	"admin.diary_entries.show",
	AdminDiaryEntryPrefix,
)
var AdminDiaryEntryNew = routing.NewSimpleRoute(
	"/new",
	"admin.diary_entries.new",
	AdminDiaryEntryPrefix,
)
var AdminDiaryEntryToday = routing.NewSimpleRoute(
	"/today",
	"admin.diary_entries.today",
	AdminDiaryEntryPrefix,
)
var AdminDiaryEntryCreate = routing.NewSimpleRoute(
	"",
	"admin.diary_entries.create",
	AdminDiaryEntryPrefix,
)
var AdminDiaryEntryEdit = routing.NewRouteWithSerialID(
	"/:id/edit",
	"admin.diary_entries.edit",
	AdminDiaryEntryPrefix,
)
var AdminDiaryEntryUpdate = routing.NewRouteWithSerialID(
	"/:id",
	"admin.diary_entries.update",
	AdminDiaryEntryPrefix,
)
var AdminDiaryEntryDestroy = routing.NewRouteWithSerialID(
	"/:id",
	"admin.diary_entries.destroy",
	AdminDiaryEntryPrefix,
)
