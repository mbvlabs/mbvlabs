// Package services provides application services for business workflows.
package services

import (
	"mbvlabs/clients/firecrawl"
	"mbvlabs/clients/serper"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"services",
	fx.Provide(
		NewIdentity,
	),
	fx.Provide(
		NewWorks,
	),
	fx.Provide(
		NewProjects,
		NewBlogPosts,
	),
	fx.Provide(
		firecrawl.New,
		serper.New,
		NewSearch,
	),
)
