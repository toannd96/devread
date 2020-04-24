package middleware

import (
	"backend-viblo-trending/model"
	"net/http"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/labstack/echo"
)

func LimitMiddleware(lmt *limiter.Limiter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			httpErr := tollbooth.LimitByRequest(lmt, c.Response(), c.Request())
			if httpErr != nil {
				return c.JSON(http.StatusTooManyRequests, model.Response{
					StatusCode: http.StatusTooManyRequests,
					Message:    "Bạn đã đạt đến giới hạn yêu cầu tối đa",
					Data:       nil,
				})
			}
			return next(c)
		}
	}
}
