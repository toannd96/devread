package handler

import (
	"backend-viblo-trending/log"
	"backend-viblo-trending/model"
	"backend-viblo-trending/model/requests"
	security "backend-viblo-trending/security"
	"net/http"

	validator "github.com/go-playground/validator/v10"
	uuid "github.com/google/uuid"
	"github.com/labstack/echo"
)

type UserHandler struct {
}

func (u *UserHandler) HandleSignUp(c echo.Context) error {

	req := requests.ReqSignUp{}
	if err := c.Bind(&req); err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	hash := security.HashAndSalt([]byte(req.Password))
	role := model.MEMBER.String()

	userId, err := uuid.NewUUID()
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusForbidden, model.Response{
			StatusCode: http.StatusForbidden,
			Message:    err.Error(),
			Data:       nil,
		})
	}

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

func (u *UserHandler) HandleSignIn(c echo.Context) error {

	return c.JSON(http.StatusOK, echo.Map{
		"user":  "Toan",
		"email": "toan@gmail.com",
	})
}
