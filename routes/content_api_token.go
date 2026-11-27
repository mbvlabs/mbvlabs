package routes

import "github.com/mbvlabs/andurel/pkg/routing"

var (
	AdminContentTokenShow   = routing.NewSimpleRoute("/content-token", "admin.content_token.show", AdminPrefix, routing.TypeScriptRoute())
	AdminContentTokenCreate = routing.NewSimpleRoute("/content-token", "admin.content_token.create", AdminPrefix, routing.TypeScriptRoute())
	AdminContentTokenDelete = routing.NewSimpleRoute("/content-token", "admin.content_token.delete", AdminPrefix, routing.TypeScriptRoute())
)
