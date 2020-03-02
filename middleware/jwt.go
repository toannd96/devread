package middleware

import (
	"backend-viblo-trending/model"
	"backend-viblo-trending/security"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func JwtMiddleware() echo.MiddlewareFunc {
	config := middleware.JWTConfig{
		Claims:     &model.JwtCustomClaims{},
		SigningKey: []byte(security.SECRET_KEY),
	}

	return middleware.JWTWithConfig(config)

}
