package routes

import (
	"mbvlabs/internal/routing"
)

const ApiArticlePrefix = "/api/articles"

var ApiArticleCreate = routing.NewSimpleRoute(
	"",
	"api.articles.create",
	ApiArticlePrefix,
)
