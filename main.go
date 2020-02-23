package main

import (
	"backend-viblo-trending/db"
	"backend-viblo-trending/handler"
	"os"

	log "backend-viblo-trending/log"

	"github.com/labstack/echo"
)

func init() {
	os.Setenv("APP_NAME", "github")
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

	log.Error("Co loi xay ra")

	e := echo.New()
	e.GET("/", handler.Welcome)
	e.GET("/user/sign-in", handler.HandleSingIn)
	e.GET("/user/sign-up", handler.HandleSingUp)
	e.Logger.Fatal(e.Start(":3000"))
}
