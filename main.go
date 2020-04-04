package main

import (
	"backend-viblo-trending/db"
	"backend-viblo-trending/handler"
	"backend-viblo-trending/helper"
	"backend-viblo-trending/repository/repo_impl"
	"backend-viblo-trending/router"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("không nhận được biến môi trường")
	}
}

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

	customValidator := helper.NewCustomValidator()
	customValidator.RegisterValidate()

	e.Validator = customValidator

	userHandler := handler.UserHandler{
		UserRepo: repo_impl.NewUserRepo(sql),
		AuthRepo: repo_impl.NewAuthRepo(client),
	}

	oauthGithub := handler.OauthGithub{}

	api := router.API{
		Echo:        e,
		UserHandler: userHandler,
		OauthGithub: oauthGithub,
	}

	api.SetupRouter()

	e.Logger.Fatal(e.Start(":3000"))
}
