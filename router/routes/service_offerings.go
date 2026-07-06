package routes

import (
	"mbvlabs/internal/routing"
)

const ServiceOfferingPrefix = "/services"
const LegacyServiceOfferingPrefix = "/service-offerings"

var ServiceOfferingIndex = routing.NewSimpleRoute(
	"",
	"service_offerings.index",
	ServiceOfferingPrefix,
)

var LegacyServiceOfferingIndex = routing.NewSimpleRoute(
	"",
	"service_offerings.index.legacy",
	LegacyServiceOfferingPrefix,
)

var LegacyServiceOfferingShow = routing.NewRouteWithSerialID(
	"/:id",
	"service_offerings.show.legacy",
	LegacyServiceOfferingPrefix,
)
