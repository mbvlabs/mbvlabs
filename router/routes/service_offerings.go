package routes

import (
	"mbvlabs/internal/routing"
)

const ServiceOfferingPrefix = "/services"

var ServiceOfferingIndex = routing.NewSimpleRoute(
	"",
	"service_offerings.index",
	ServiceOfferingPrefix,
)
