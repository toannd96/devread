package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func BodyLimitMiddleware() echo.MiddlewareFunc {
	config := middleware.BodyLimitConfig{
		Limit:   "2M",
	}
	return middleware.BodyLimitWithConfig(config)
}
