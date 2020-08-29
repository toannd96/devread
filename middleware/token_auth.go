package middleware

import (
	"tech_posts_trending/model"
	"tech_posts_trending/security"
	"github.com/labstack/echo/v4"
	"net/http"
)

func TokenAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := security.AccessTokenValid(c.Request())
			if err != nil {
				return c.JSON(http.StatusUnauthorized, model.Response{
					StatusCode: http.StatusUnauthorized,
					Message:    "Truy cập không được phép",
				})
			}
			return next(c)
		}
	}
}
