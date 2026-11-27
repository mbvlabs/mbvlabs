package routes

import (
	"github.com/mbvlabs/andurel/pkg/routing"
)

const ProjectInquiryPrefix = "/contact"
const LegacyProjectInquiryPrefix = "/project-inquiries"

var ProjectInquiryIndex = routing.NewSimpleRoute(
	"",
	"project_inquiries.index",
	ProjectInquiryPrefix,
	routing.TypeScriptRoute(),
)
var ProjectInquiryCreate = routing.NewSimpleRoute(
	"",
	"project_inquiries.create",
	ProjectInquiryPrefix,
	routing.TypeScriptRoute(),
)

var LegacyProjectInquiryIndex = routing.NewSimpleRoute(
	"",
	"project_inquiries.index.legacy",
	LegacyProjectInquiryPrefix,
)

var LegacyProjectInquiryCreate = routing.NewSimpleRoute(
	"",
	"project_inquiries.create.legacy",
	LegacyProjectInquiryPrefix,
)
