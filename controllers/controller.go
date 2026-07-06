// Package controllers provides HTTP handlers for the web application.
package controllers

import (
	"mbvlabs/controllers/admin"
	"mbvlabs/router"

	"go.uber.org/fx"
)

var otherCache = NewCacheBuilder[string]().WithSize(2).Build

var constructors = fx.Provide(
	otherCache,
	NewPages,
	NewAssets,
	NewAPI,
	NewSessions,
	NewConfirmations,
	NewResetPasswords,
	admin.NewProjectInquiries,
	admin.NewBlogPosts,
	NewBlogPosts,
	admin.NewProjects,
	NewProjects,
	NewProjectInquiries,
	admin.NewWorks,
	NewWorks,
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
)
