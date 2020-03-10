package router

import (
	"backend-viblo-trending/handler"
	middleware "backend-viblo-trending/middleware"

	"github.com/labstack/echo"
)

type API struct {
	Echo        *echo.Echo
	UserHandler handler.UserHandler
}

func (api *API) SetupRouter() {

	// user
	api.Echo.POST("/user/sign-in", api.UserHandler.SignIn)
	api.Echo.POST("/user/sign-up", api.UserHandler.SignUp)
	api.Echo.POST("user/sign-out", api.UserHandler.SignOut, middleware.JwtAccessToken(), middleware.JwtRefreshToken())
	api.Echo.GET("/user/refresh-token", api.UserHandler.RefeshToken, middleware.JwtRefreshToken())

	// profile
	user := api.Echo.Group("/user", middleware.JwtAccessToken())
	user.GET("/profile", api.UserHandler.Profile)
	user.PUT("/profile/update", api.UserHandler.UpdateProfile)

}
