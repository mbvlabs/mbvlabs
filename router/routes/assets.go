// Package routes contains all the routes used throughout the project
package routes

import (
	"fmt"
	"time"

	"mbvlabs/internal/routing"
)

const AssetsPrefix = "assets"

var startTime = time.Now().Unix()

var Robots = routing.NewSimpleRoute(
	"/robots.txt",
	"robots",
	"",
)

var Sitemap = routing.NewSimpleRoute(
	"/sitemap.xml",
	"sitemap",
	"",
)

var Stylesheet = routing.NewSimpleRoute(
	fmt.Sprintf("/css/%v/style.css", startTime),
	"css.stylesheet",
	AssetsPrefix,
)

var Scripts = routing.NewSimpleRoute(
	fmt.Sprintf("/js/%v/scripts.js", startTime),
	"js.scripts",
	AssetsPrefix,
)

var Script = routing.NewRouteWithFile(
	fmt.Sprintf("/js/%v/:file", startTime),
	"js.script",
	AssetsPrefix,
)
