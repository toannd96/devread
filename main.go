package main

import (
	"backend-viblo-trending/db"
	"backend-viblo-trending/log"
	_ "backend-viblo-trending/docs"
	"backend-viblo-trending/handler"
	"backend-viblo-trending/helper"
	"backend-viblo-trending/repository/repo_impl"
	"backend-viblo-trending/router"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func init() {
	fmt.Println("PRODUCTION ENVIROMENT")
	os.Setenv("APP_NAME", "post-trending")
	log.InitLogger(false)

	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("không nhận được biến môi trường")
	}
}

// @title Github Trending API
// @version 1.0
// @description Secure REST API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey token-verify-account
// @in query
// @name token

// @securityDefinitions.apikey token-reset-password
// @in query
// @name token

// @host localhost:3000
// @BasePath /

func main() {

	// redis details
	redisHost     := "host.docker.internal"
	redisPort := os.Getenv("REDIS_PORT")

	// postgres details
	host     := "host.docker.internal"
	port     := os.Getenv("DB_PORT")
	password := os.Getenv("DB_PASSWORD")
	username := os.Getenv("DB_USERNAME")
	dbname   := os.Getenv("DB_NAME")

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

	repoHandler := handler.RepoHandler{
		GithubRepo: repo_impl.NewGithubRepo(sql),
		AuthRepo:   repo_impl.NewAuthRepo(client),
	}

	api := router.API{
		Echo:        e,
		UserHandler: userHandler,
		RepoHandler: repoHandler,
	}

	api.SetupRouter()

	go scheduleUpdateTrending(60*time.Second, repoHandler)

	e.Logger.Fatal(e.Start(":3000"))
}

func scheduleUpdateTrending(timeSchedule time.Duration, handler handler.RepoHandler) {
	ticker := time.NewTicker(timeSchedule)
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("Checking from github...")
				helper.CrawlRepo(handler.GithubRepo)
			}
		}
	}()
}
