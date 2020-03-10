package middleware

import (
	"backend-viblo-trending/model"
	"backend-viblo-trending/model/requests"
	"net/http"

	"github.com/labstack/echo"
)

func IsAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := requests.RequestSignIn{}
			if err := c.Bind(&req); err != nil {
				return c.JSON(http.StatusBadRequest, model.Response{
					StatusCode: http.StatusBadRequest,
					Message:    err.Error(),
					Data:       nil,
				})
			}

			if req.Email != "admin@gmail.com" {
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
