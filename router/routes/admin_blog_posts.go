package routes

import (
	"mbvlabs/internal/routing"
)

const AdminBlogPostPrefix = "/admin/blog-posts"

var AdminBlogPostIndex = routing.NewSimpleRoute(
	"",
	"admin.blog_posts.index",
	AdminBlogPostPrefix,
)
var AdminBlogPostShow = routing.NewRouteWithSerialID(
	"/:id",
	"admin.blog_posts.show",
	AdminBlogPostPrefix,
)
var AdminBlogPostNew = routing.NewSimpleRoute(
	"/new",
	"admin.blog_posts.new",
	AdminBlogPostPrefix,
)
var AdminBlogPostCreate = routing.NewSimpleRoute(
	"",
	"admin.blog_posts.create",
	AdminBlogPostPrefix,
)
var AdminBlogPostEdit = routing.NewRouteWithSerialID(
	"/:id/edit",
	"admin.blog_posts.edit",
	AdminBlogPostPrefix,
)
var AdminBlogPostUpdate = routing.NewRouteWithSerialID(
	"/:id",
	"admin.blog_posts.update",
	AdminBlogPostPrefix,
)
var AdminBlogPostDestroy = routing.NewRouteWithSerialID(
	"/:id",
	"admin.blog_posts.destroy",
	AdminBlogPostPrefix,
)
