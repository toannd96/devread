package middleware

import "github.com/labstack/echo"

func IsAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// handle logic
			if 1 == 1 {
				// gia su co loi
				return c.JSON(500, echo.Map{
					"err": "Bạn chưa đăng nhập",
				})
			}
			return next(c)
		}
	}

}
