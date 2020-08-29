package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func GzipMiddleware() echo.MiddlewareFunc {
	config := middleware.GzipConfig{
		Level: 5,
	}
	return middleware.GzipWithConfig(config)
}
