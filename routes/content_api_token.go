package routes

import "github.com/mbvlabs/andurel/pkg/routing"

var (
	AdminContentTokenShow   = routing.NewSimpleRoute("/content-token", "admin.content_token.show", AdminPrefix)
	AdminContentTokenCreate = routing.NewSimpleRoute("/content-token", "admin.content_token.create", AdminPrefix)
	AdminContentTokenDelete = routing.NewSimpleRoute("/content-token", "admin.content_token.delete", AdminPrefix)
)
