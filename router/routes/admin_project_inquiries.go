package routes

import (
	"mbvlabs/internal/routing"
)

const AdminProjectInquiryPrefix = "/admin/project-inquiries"

var AdminProjectInquiryIndex = routing.NewSimpleRoute(
	"",
	"admin.project_inquiries.index",
	AdminProjectInquiryPrefix,
)
var AdminProjectInquiryShow = routing.NewRouteWithSerialID(
	"/:id",
	"admin.project_inquiries.show",
	AdminProjectInquiryPrefix,
)
var AdminProjectInquiryNew = routing.NewSimpleRoute(
	"/new",
	"admin.project_inquiries.new",
	AdminProjectInquiryPrefix,
)
var AdminProjectInquiryCreate = routing.NewSimpleRoute(
	"",
	"admin.project_inquiries.create",
	AdminProjectInquiryPrefix,
)
var AdminProjectInquiryEdit = routing.NewRouteWithSerialID(
	"/:id/edit",
	"admin.project_inquiries.edit",
	AdminProjectInquiryPrefix,
)
var AdminProjectInquiryUpdate = routing.NewRouteWithSerialID(
	"/:id",
	"admin.project_inquiries.update",
	AdminProjectInquiryPrefix,
)
var AdminProjectInquiryDestroy = routing.NewRouteWithSerialID(
	"/:id",
	"admin.project_inquiries.destroy",
	AdminProjectInquiryPrefix,
)
