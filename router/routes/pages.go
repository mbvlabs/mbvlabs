package routes

import (
	"mbvlabs/internal/routing"
)

var HomePage = routing.NewSimpleRoute(
	"/",
	"pages.home",
	"",
)
var AboutMe = routing.NewSimpleRoute(
	"/about",
	"pages.about_me",
	"",
)

var LegacyAboutMe = routing.NewSimpleRoute(
	"/about-me",
	"pages.about_me.legacy",
	"",
)
