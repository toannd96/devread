package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func CORSMiddleware() echo.MiddlewareFunc {
	config := middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderContentLength, echo.HeaderAccept},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}
	return middleware.CORSWithConfig(config)
}
