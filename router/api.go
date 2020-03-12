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
	api.Echo.POST("/user/refresh-token", api.UserHandler.RefeshToken, middleware.RenewToken())
	// api.Echo.POST("/user/sign-out", api.UserHandler.SignOut, middleware.IsLoggedIn(), middleware.RenewToken())

	// profile
	user := api.Echo.Group("/user", middleware.IsLoggedIn())
	user.GET("/profile", api.UserHandler.Profile)
	user.PUT("/profile/update", api.UserHandler.UpdateProfile)
}
