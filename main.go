package main

import (
	"devread/crawler"
	"devread/db"
	_ "devread/docs"
	"devread/handler"
	"devread/helper"
	"devread/log"
	"devread/repository/repo_impl"
	"devread/router"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("không nhận được biến môi trường")
	}
	log.InitLogger(false)
}

// @title DevRead API
// @version 1.0
// @description Nền tảng tổng hợp kiến thức cho developer
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey jwt
// @in header
// @name Authorization

// @host localhost:3000
// @BasePath /

func main() {

	// redis details
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	// postgres details
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	password := os.Getenv("DB_PASSWORD")
	username := os.Getenv("DB_USERNAME")
	dbname := os.Getenv("DB_NAME")

	// connect redis
	client := &db.RedisDB{
		Host: redisHost,
		Port: redisPort,
	}
	client.NewRedisDB()

	// connect postgres
	sql := &db.Sql{
		Host:     host,
		Port:     port,
		UserName: username,
		Password: password,
		DbName:   dbname,
	}
	sql.Connect()
	defer sql.Close()

	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	customValidator := helper.NewCustomValidator()
	customValidator.RegisterValidate()

	e.Validator = customValidator

	userHandler := handler.UserHandler{
		UserRepo: repo_impl.NewUserRepo(sql),
		AuthRepo: repo_impl.NewAuthenRepo(client),
	}

	postHandler := handler.PostHandler{
		PostRepo: repo_impl.NewPostRepo(sql),
		AuthRepo: repo_impl.NewAuthenRepo(client),
	}

	api := router.API{
		Echo:        e,
		UserHandler: userHandler,
		PostHandler: postHandler,
	}

	api.SetupRouter()

	// time start crawler
	go crawler.VibloPost(postHandler.PostRepo)
	go crawler.CodeaholicguyPost(postHandler.PostRepo)
	go crawler.QuancamPost(postHandler.PostRepo)
	go crawler.ThefullsnackPost(postHandler.PostRepo)
	go crawler.ToidicodedaoPost(postHandler.PostRepo)
	crawler.YellowcodePost(postHandler.PostRepo)

	// schedule crawler
	go schedule(30*time.Minute, postHandler, 1)
	go schedule(1*time.Hour, postHandler, 2)
	go schedule(24*time.Hour, postHandler, 3)
	go schedule(24*time.Hour, postHandler, 4)
	go schedule(24*time.Hour, postHandler, 5)
	schedule(24*time.Hour, postHandler, 6)

	e.Logger.Fatal(e.Start(":3000"))
}

func schedule(timeSchedule time.Duration, handler handler.PostHandler, crowIlnndex int) {
	ticker := time.NewTicker(timeSchedule)
	func() {
		for {
			switch crowIlnndex {
			case 1:
				<-ticker.C
				crawler.VibloPost(handler.PostRepo)
			case 2:
				<-ticker.C
				crawler.ToidicodedaoPost(handler.PostRepo)
			case 3:
				<-ticker.C
				crawler.ThefullsnackPost(handler.PostRepo)
			case 4:
				<-ticker.C
				crawler.QuancamPost(handler.PostRepo)
			case 5:
				<-ticker.C
				crawler.CodeaholicguyPost(handler.PostRepo)
			case 6:
				<-ticker.C
				crawler.YellowcodePost(handler.PostRepo)
			}
		}
	}()
}
