package middleware

import (
	"backend-viblo-trending/model"
	"backend-viblo-trending/security"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func IsLoggedIn() echo.MiddlewareFunc {
	config := middleware.JWTConfig{
		Claims:      &model.JwtCustomClaims{},
		SigningKey:  []byte(security.SECRET_KEY),
		TokenLookup: "cookie:AccessToken",
		BeforeFunc: func(c echo.Context) {
			accessToken, _ := c.Cookie("AccessToken")
			if accessToken == nil {
				refreshToken, _ := c.Cookie("RefreshToken")
				if refreshToken == nil {
					return
				}
			}
		},
	}
	return middleware.JWTWithConfig(config)
}

func RenewToken() echo.MiddlewareFunc {
	config := middleware.JWTConfig{
		Claims:      &model.JwtCustomClaims{},
		SigningKey:  []byte(security.SECRET_KEY),
		TokenLookup: "cookie:RefreshToken",
		BeforeFunc: func(c echo.Context) {
			refreshToken, _ := c.Cookie("RefreshToken")
			if refreshToken == nil {
				return
			}
		},
	}
	return middleware.JWTWithConfig(config)
}
