package main

import (
	"backend-viblo-trending/db"
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

	e.Logger.Fatal(e.Start(":3000"))
}
