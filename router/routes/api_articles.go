package routes

import (
	"mbvlabs/internal/routing"
)

const ApiArticlePrefix = "/api/articles"

var (
	ApiArticleIndex = routing.NewSimpleRoute(
		"",
		"api.articles.index",
		ApiArticlePrefix,
	)
	ApiArticleCreate = routing.NewSimpleRoute(
		"",
		"api.articles.create",
		ApiArticlePrefix,
	)
)
