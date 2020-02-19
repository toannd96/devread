package handler

import (
	"net/http"

	"github.com/labstack/echo"
)

func HandleSingIn(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"user":  "Toan",
		"email": "toan@gmail.com",
	})
}

func HandleSingUp(c echo.Context) error {
	type User struct {
		Email    string
		FullName string
	}
	user := User{
		Email:    "toan@gmail.com",
		FullName: "Toan",
	}
	return c.JSON(http.StatusOK, user)
}
