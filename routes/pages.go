package routes

import (
	"github.com/mbvlabs/andurel/pkg/routing"
)

var HomePage = routing.NewSimpleRoute(
	"/",
	"pages.home",
	"",
	routing.TypeScriptRoute(),
)

var AboutMe = routing.NewSimpleRoute("/about", "pages.about_me", "", routing.TypeScriptRoute())
var LegacyAboutMe = routing.NewSimpleRoute("/about-me", "pages.about_me.legacy", "")
