package routes

import (
	"github.com/mbvlabs/andurel/pkg/routing"
)

const ServiceOfferingPrefix = "/services"

var ServiceOfferingIndex = routing.NewSimpleRoute(
	"",
	"service_offerings.index",
	ServiceOfferingPrefix,
	routing.TypeScriptRoute(),
)
