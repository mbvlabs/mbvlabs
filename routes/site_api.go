package routes

import "github.com/mbvlabs/andurel/pkg/routing"

var (
	ApiArticleIndex              = routing.NewSimpleRoute("/articles", "api.articles.index", APIPrefix)
	ApiArticleCreate             = routing.NewSimpleRoute("/articles", "api.articles.create", APIPrefix)
	ApiWorkCreate                = routing.NewSimpleRoute("/works", "api.works.create", APIPrefix)
	ApiDiaryThoughtsCurrentWeek  = routing.NewSimpleRoute("/diary/thoughts/current-week", "api.diary.thoughts.current_week", APIPrefix)
	ApiDiaryThoughtsPreviousWeek = routing.NewSimpleRoute("/diary/thoughts/previous-week", "api.diary.thoughts.previous_week", APIPrefix)
	ApiSearch                    = routing.NewSimpleRoute("/search", "api.search", APIPrefix)
	ApiScrape                    = routing.NewSimpleRoute("/scrape", "api.scrape", APIPrefix)
	ApiCrawl                     = routing.NewSimpleRoute("/crawl", "api.crawl.create", APIPrefix)
	ApiCrawlStatus               = routing.NewRouteWithStringID("/crawl/:id", "api.crawl.show", APIPrefix)
	ApiMap                       = routing.NewSimpleRoute("/map", "api.map", APIPrefix)
)
