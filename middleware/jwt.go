package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"os"

	"devread/model"
)

func JWTMiddleware() echo.MiddlewareFunc {
	config := middleware.JWTConfig{
		Claims:     &model.TokenDetails{},
		SigningKey: []byte(os.Getenv("ACCESS_SECRET")),
	}

	return middleware.JWTWithConfig(config)
}
