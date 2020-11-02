package main

import (
	"fmt"
	"os"
	"tech_posts_trending/db"
	_ "tech_posts_trending/docs"
	"tech_posts_trending/handler"
	"tech_posts_trending/helper"
	"tech_posts_trending/log"
	"tech_posts_trending/repository/repo_impl"
	"tech_posts_trending/router"
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

// @title Tech Posts Trending API
// @version 1.0
// @description Build web applications to crawler article information on technology blogs using Echo Framework (Golang)
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

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
		AuthRepo: repo_impl.NewAuthRepo(client),
	}

	postHandler := handler.PostHandler{
		PostRepo: repo_impl.NewPostRepo(sql),
		AuthRepo: repo_impl.NewAuthRepo(client),
	}

	api := router.API{
		Echo:        e,
		UserHandler: userHandler,
		PostHandler: postHandler,
	}

	api.SetupRouter()

	go scheduleUpdateTrending(24*time.Second, postHandler)

	e.Logger.Fatal(e.Start(":3000"))
}

func scheduleUpdateTrending(timeSchedule time.Duration, handler handler.PostHandler) {
	ticker := time.NewTicker(timeSchedule)
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("Quét bài viết ...")
				helper.VibloPost(handler.PostRepo)
				helper.ToidicodedaoPost(handler.PostRepo)
				helper.ThefullsnackPost(handler.PostRepo)
				helper.QuancamPost(handler.PostRepo)
				helper.CodeaholicguyPost(handler.PostRepo)
				helper.YellowcodePost(handler.PostRepo)
			}
		}
	}()
}
