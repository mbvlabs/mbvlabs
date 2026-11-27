package routes

import (
	"github.com/mbvlabs/andurel/pkg/routing"
)

const BlogPostPrefix = "/blog"
const LegacyBlogPostPrefix = "/blog-posts"

var BlogPostIndex = routing.NewSimpleRoute(
	"",
	"blog_posts.index",
	BlogPostPrefix,
	routing.TypeScriptRoute(),
)

var BlogPostFeed = routing.NewSimpleRoute(
	"/rss.xml",
	"blog_posts.feed",
	"",
	routing.TypeScriptRoute(),
)

var BlogPostShow = routing.NewRouteWithSlug(
	"/:slug",
	"blog_posts.show",
	BlogPostPrefix,
	routing.TypeScriptRoute(),
)

var LegacyBlogPostIndex = routing.NewSimpleRoute(
	"",
	"blog_posts.index.legacy",
	LegacyBlogPostPrefix,
)

var LegacyBlogPostShow = routing.NewRouteWithSlug(
	"/:slug",
	"blog_posts.show.legacy",
	LegacyBlogPostPrefix,
)
