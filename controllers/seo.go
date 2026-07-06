package controllers

import "github.com/labstack/echo/v5"

func setPrivateSEOHeaders(etx *echo.Context) {
	etx.Response().Header().Set("X-Robots-Tag", "noindex, nofollow")
}
