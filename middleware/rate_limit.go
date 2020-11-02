package middleware

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
	"net/http"
	"tech_posts_trending/model"
)

var limiter = rate.NewLimiter(1, 1)

func LimitRequest() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if limiter.Allow() == false {
				return c.JSON(http.StatusTooManyRequests, model.Response{
					StatusCode: http.StatusTooManyRequests,
					Message:    "Đã đạt đến giới hạn yêu cầu tối đa 1 request/s",
				})
			}
			return next(c)
		}
	}
}
