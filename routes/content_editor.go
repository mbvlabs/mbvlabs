package routes

import "github.com/mbvlabs/andurel/pkg/routing"

var (
	ContentArticles = routing.NewSimpleRoute("/content/articles", "api.content.articles.index", APIPrefix)
	ContentArticle  = routing.NewRouteWithSlug("/content/articles/:slug", "api.content.articles.show", APIPrefix)
	ContentProjects = routing.NewSimpleRoute("/content/projects", "api.content.projects.index", APIPrefix)
	ContentProject  = routing.NewRouteWithSlug("/content/projects/:slug", "api.content.projects.show", APIPrefix)
	ContentWorks    = routing.NewSimpleRoute("/content/works", "api.content.works.index", APIPrefix)
	ContentWork     = routing.NewRouteWithSlug("/content/works/:slug", "api.content.works.show", APIPrefix)
)
