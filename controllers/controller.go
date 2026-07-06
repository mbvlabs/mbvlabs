// Package controllers provides HTTP handlers for the web application.
package controllers

import (
	"mbvlabs/controllers/admin"
	"mbvlabs/controllers/api"
	"mbvlabs/router"
	"time"

	"go.uber.org/fx"
)

var otherCache = NewCacheBuilder[string]().
	WithSize(2).
	WithDefaultTTL(72 * time.Hour).
	Build

func newSitemapCacheInvalidator(cache *Cache[string]) admin.SitemapCacheInvalidator {
	return func() {
		cache.Invalidate(sitemapCacheKey)
	}
}

var constructors = fx.Provide(
	otherCache,
	newSitemapCacheInvalidator,
	NewPages,
	NewAssets,
	NewAPI,
	NewSessions,
	NewConfirmations,
	NewResetPasswords,
	admin.NewAdmin,
	admin.NewProjectInquiries,
	admin.NewBlogPosts,
	NewBlogPosts,
	admin.NewProjects,
	NewProjects,
	NewProjectInquiries,
	admin.NewWorks,
	NewWorks,
	api.NewWorks,
	NewServiceOfferings,
)

var Module = fx.Module(
	"controllers",
	constructors,
	fx.Invoke(func(r *router.Router, c Pages) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c Assets) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c API) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c Sessions) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c Confirmations) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c ResetPasswords) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c admin.Admin) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c admin.ProjectInquiries) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c ProjectInquiries) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c admin.BlogPosts) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c BlogPosts) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c admin.Projects) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c Projects) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c admin.Works) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c Works) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c api.Works) error {
		return c.RegisterRoutes(r)
	}),
	fx.Invoke(func(r *router.Router, c ServiceOfferings) error {
		return c.RegisterRoutes(r)
	}),
)
