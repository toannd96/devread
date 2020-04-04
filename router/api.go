package router

import (
	"backend-viblo-trending/handler"
	"backend-viblo-trending/middleware"

	"github.com/didip/tollbooth"
	"github.com/labstack/echo"
)

type API struct {
	Echo        *echo.Echo
	UserHandler handler.UserHandler
	OauthGithub handler.OauthGithub
}

func (api *API) SetupRouter() {

	limit := tollbooth.NewLimiter(1, nil)

	// user
	user := api.Echo.Group("/user", middleware.CORSMiddleware(), middleware.LimitMiddleware(limit))
	user.POST("/sign-in", api.UserHandler.SignIn)
	user.POST("/sign-up", api.UserHandler.SignUp)
	user.POST("/refresh", api.UserHandler.Refresh)
	user.GET("/github/sign-in", api.OauthGithub.GithubSignIn)
	user.GET("/github/callback", api.OauthGithub.GithubCallback)

	// user profile
	user_profile := api.Echo.Group("/user", middleware.CORSMiddleware(), middleware.TokenAuthMiddleware(), middleware.LimitMiddleware(limit))
	user_profile.GET("/profile", api.UserHandler.Profile)
	user_profile.PUT("/profile/update", api.UserHandler.UpdateProfile)
	user_profile.POST("/sign-out", api.UserHandler.SignOut)
}
