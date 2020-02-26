package main

import (
	"backend-viblo-trending/db"
	"backend-viblo-trending/handler"
	"backend-viblo-trending/repository/repo_impl"
	"backend-viblo-trending/router"
	"os"

	log "backend-viblo-trending/log"

	"github.com/labstack/echo"
)

func init() {
	os.Setenv("APP_NAME", "viblo")
	log.InitLogger(false)
}

func main() {

	sql := &db.Sql{
		Host:     "localhost",
		Port:     5432,
		UserName: "postgres",
		Password: "nguyendactoan",
		DbName:   "golang",
	}
	sql.Connect()
	defer sql.Close()

	e := echo.New()

	userHandler := handler.UserHandler{
		UserRepo: repo_impl.NewUserRepo(sql),
	}
	api := router.API{
		Echo:        e,
		UserHandler: userHandler,
	}

	api.SetupRouter()

	e.Logger.Fatal(e.Start(":3000"))
}
