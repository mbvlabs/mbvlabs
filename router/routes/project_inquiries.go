package routes

import (
	"mbvlabs/internal/routing"
)

const ProjectInquiryPrefix = "/contact"
const LegacyProjectInquiryPrefix = "/project-inquiries"

var ProjectInquiryIndex = routing.NewSimpleRoute(
	"",
	"project_inquiries.index",
	ProjectInquiryPrefix,
)
var ProjectInquiryCreate = routing.NewSimpleRoute(
	"",
	"project_inquiries.create",
	ProjectInquiryPrefix,
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
