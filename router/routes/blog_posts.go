package routes

import (
	"mbvlabs/internal/routing"
)

const BlogPostPrefix = "/blog"
const LegacyBlogPostPrefix = "/blog-posts"

var BlogPostIndex = routing.NewSimpleRoute(
	"",
	"blog_posts.index",
	BlogPostPrefix,
)

var BlogPostFeed = routing.NewSimpleRoute(
	"/rss.xml",
	"blog_posts.feed",
	"",
)

var BlogPostShow = routing.NewRouteWithSlug(
	"/:slug",
	"blog_posts.show",
	BlogPostPrefix,
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
