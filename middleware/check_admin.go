package middleware

import (
	"backend-viblo-trending/model"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func IsAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			isAdmin := claims["Role"].(string)
			if isAdmin != "ADMIN" {
				return c.JSON(http.StatusBadRequest, model.Response{
					StatusCode: http.StatusBadRequest,
					Message:    "Bạn không có quyền gọi API này",
					Data:       nil,
				})
			}
			return next(c)
		}
	}

}
