package router

import (
	"tech_posts_trending/handler"
	"tech_posts_trending/middleware"
	"github.com/labstack/echo/v4"
)

type API struct {
	Echo        *echo.Echo
	UserHandler handler.UserHandler
	PostHandler handler.PostHandler
}

func (api *API) SetupRouter() {

	// user
	user := api.Echo.Group("/user",
		middleware.CORSMiddleware(),
		middleware.LimitRequest(),
		middleware.HeadersMiddleware(),
		middleware.HeadersAccept(),
		middleware.GzipMiddleware(),
		)
	user.POST("/sign-in", api.UserHandler.SignIn)
	user.POST("/sign-up", api.UserHandler.SignUp)
	user.POST("/refresh", api.UserHandler.Refresh)
	user.POST("/verify", api.UserHandler.VerifyAccount)
	user.POST("/password/forgot", api.UserHandler.ForgotPassword)
	user.PUT("/password/reset", api.UserHandler.ResetPassword)

	// user profile
	userProfile := api.Echo.Group("/user",
		middleware.CORSMiddleware(),
		middleware.TokenAuthMiddleware(),
		middleware.LimitRequest(),
		middleware.HeadersMiddleware(),
		middleware.HeadersAccept(),
		middleware.GzipMiddleware(),
		)
	userProfile.GET("/profile", api.UserHandler.Profile)
	userProfile.PUT("/profile/update", api.UserHandler.UpdateProfile)
	userProfile.POST("/sign-out", api.UserHandler.SignOut)

	// bookmark user
	bookmark := api.Echo.Group("/user",
		middleware.CORSMiddleware(),
		middleware.TokenAuthMiddleware(),
		middleware.LimitRequest(),
		middleware.HeadersMiddleware(),
		middleware.HeadersAccept(),
		middleware.GzipMiddleware(),
		)
	bookmark.GET("/bookmark/list", api.PostHandler.SelectBookmarks)
	bookmark.POST("/bookmark/add", api.PostHandler.Bookmark)
	bookmark.DELETE("/bookmark/delete", api.PostHandler.DelBookmark)

	// post
	post := api.Echo.Group("/",
		middleware.CORSMiddleware(),
		middleware.LimitRequest(),
		middleware.HeadersMiddleware(),
		middleware.HeadersAccept(),
		middleware.GzipMiddleware(),
	)
	post.GET("posts", api.PostHandler.PostTrending)
}
