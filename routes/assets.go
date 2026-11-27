package routes

import (
	"github.com/mbvlabs/andurel/pkg/routing"
)

const AssetsPrefix = "/assets"

var Robots = routing.NewSimpleRoute(
	"/robots.txt",
	"assets.robots",
	"",
)

var Sitemap = routing.NewSimpleRoute(
	"/sitemap.xml",
	"assets.sitemap",
	"",
)

// ViteBuild serves the Vite build output. Vite fingerprints file names, so
// the path needs no cache-busting segment.
var ViteBuild = routing.NewSimpleRoute(
	"/dist/*",
	"vite.build",
	AssetsPrefix,
)

const IndexNowKeyValue = "d12s7upn4vcbgvgeskaukp2k2h6snfmb"

var IndexNowKey = routing.NewSimpleRoute("/d12s7upn4vcbgvgeskaukp2k2h6snfmb.txt", "assets.index_now_key", "")
