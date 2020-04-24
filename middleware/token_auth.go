package middleware

import (
	"backend-viblo-trending/model"
	"backend-viblo-trending/security"
	"net/http"

	"github.com/labstack/echo"
)

func TokenAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := security.AccessTokenValid(c.Request())
			if err != nil {
				return c.JSON(http.StatusUnauthorized, model.Response{
					StatusCode: http.StatusUnauthorized,
					Message:    "Truy cập không được phép",
					Data:       nil,
				})
			}
			return next(c)
		}
	}
}
