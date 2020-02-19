package main

import (
	"backend-viblo-trending/db"
	"backend-viblo-trending/handler"

	"github.com/labstack/echo"
)

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
	e.GET("/", handler.Welcome)
	e.GET("/user/sign-in", handler.HandleSingIn)
	e.GET("/user/sign-up", handler.HandleSingUp)
	e.Logger.Fatal(e.Start(":3000"))
}
