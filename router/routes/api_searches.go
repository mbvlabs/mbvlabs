package routes

import "mbvlabs/internal/routing"

var (
	ApiSearch = routing.NewSimpleRoute(
		"/search",
		"api.search",
		APIPrefix,
	)
	ApiScrape = routing.NewSimpleRoute(
		"/scrape",
		"api.scrape",
		APIPrefix,
	)
	ApiCrawl = routing.NewSimpleRoute(
		"/crawl",
		"api.crawl.create",
		APIPrefix,
	)
	ApiCrawlStatus = routing.NewRouteWithStringID(
		"/crawl/:id",
		"api.crawl.show",
		APIPrefix,
	)
	ApiMap = routing.NewSimpleRoute(
		"/map",
		"api.map",
		APIPrefix,
	)
)
