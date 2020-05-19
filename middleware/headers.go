package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func HeadersMiddleware() echo.MiddlewareFunc {
	config := middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
	}
	return middleware.SecureWithConfig(config)
}
