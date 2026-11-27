package routes

import (
	"github.com/mbvlabs/andurel/pkg/routing"
)

const UserPrefix = "/users"

var SessionNew = routing.NewSimpleRoute(
	"/sign-in",
	"users.new_user_session",
	UserPrefix,
	routing.TypeScriptRoute(),
)

var SessionCreate = routing.NewSimpleRoute(
	"/sign-in",
	"users.user_session",
	UserPrefix,
	routing.TypeScriptRoute(),
)

// SessionDestroy signs the user out. It is registered for DELETE (the Inertia
// link) and POST (the no-JavaScript form fallback).
var SessionDestroy = routing.NewSimpleRoute(
	"/sign_out",
	"users.destroy_user_session",
	UserPrefix,
	routing.TypeScriptRoute(),
)

var PasswordNew = routing.NewSimpleRoute(
	"/password/new",
	"users.new_user_password",
	UserPrefix,
	routing.TypeScriptRoute(),
)

var PasswordCreate = routing.NewSimpleRoute(
	"/password",
	"users.user_password",
	UserPrefix,
	routing.TypeScriptRoute(),
)

var PasswordEdit = routing.NewRouteWithToken(
	"/password/:token/edit",
	"users.edit_user_password",
	UserPrefix,
	routing.TypeScriptRoute(),
)

var PasswordUpdate = routing.NewSimpleRoute(
	"/password",
	"users.update_user_password",
	UserPrefix,
	routing.TypeScriptRoute(),
)
