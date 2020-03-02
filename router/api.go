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
	api.Echo.POST("/user/sign-in", api.UserHandler.HandleSignIn)
	api.Echo.POST("/user/sign-up", api.UserHandler.HandleSignUp)

	api.Echo.GET("/user/profile", api.UserHandler.HandleProfile, middleware.JwtMiddleware())
}
