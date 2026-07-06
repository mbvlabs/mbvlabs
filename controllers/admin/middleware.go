package admin

import (
	"mbvlabs/router/middleware"

	"github.com/labstack/echo/v5"
)

var authOnly = []echo.MiddlewareFunc{middleware.AuthOnly}
