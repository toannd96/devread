package router

import (
	"backend-viblo-trending/handler"

	"github.com/labstack/echo"
)

type API struct {
	Echo        *echo.Echo
	UserHandler handler.UserHandler
}

func (api *API) SetupRouter() {
	// user
	api.Echo.GET("/user/sign-in", api.UserHandler.HandleSignIn)
	api.Echo.GET("/user/sign-up", api.UserHandler.HandleSignUp)
}
