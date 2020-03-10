package middleware

import (
	"backend-viblo-trending/model"
	"backend-viblo-trending/security"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func JwtAccessToken() echo.MiddlewareFunc {
	config := middleware.JWTConfig{
		Claims:      &model.JwtCustomClaims{},
		SigningKey:  []byte(security.SECRET_KEY),
		TokenLookup: "cookie:AccessToken",
	}

	return middleware.JWTWithConfig(config)
}

func JwtRefreshToken() echo.MiddlewareFunc {
	config := middleware.JWTConfig{
		Claims:      &model.JwtCustomClaims{},
		SigningKey:  []byte(security.SECRET_KEY),
		TokenLookup: "cookie:RefreshToken",
	}

	return middleware.JWTWithConfig(config)
}
