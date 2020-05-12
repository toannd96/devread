package middleware

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func BodyLimitMiddleware() echo.MiddlewareFunc {
	config := middleware.BodyLimitConfig{
		Limit:   "2M",
	}
	return middleware.BodyLimitWithConfig(config)
}
