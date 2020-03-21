package handler

import (
	"backend-viblo-trending/custom_error"
	"backend-viblo-trending/model"
	"backend-viblo-trending/model/req"
	"backend-viblo-trending/repository"
	security "backend-viblo-trending/security"
	"log"
	"net/http"
	"time"

	uuid "github.com/google/uuid"
	"github.com/labstack/echo"
)

type UserHandler struct {
	UserRepo repository.UserRepo
	AuthRepo repository.AuthRepo
}

func (u *UserHandler) SignUp(c echo.Context) error {
	request := req.ReqtSignUp{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	hash := security.HashAndSalt([]byte(request.Password))
	role := model.MEMBER.String()

	userID, err := uuid.NewUUID()
	if err != nil {
		return c.JSON(http.StatusForbidden, model.Response{
			StatusCode: http.StatusForbidden,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	user := model.User{
		UserID:   userID.String(),
		FullName: request.FullName,
		Email:    request.Email,
		Password: hash,
		Role:     role,
	}

	user, err = u.UserRepo.SaveUser(c.Request().Context(), user)
	if err != nil {
		return c.JSON(http.StatusConflict, model.Response{
			StatusCode: http.StatusConflict,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Đăng ký thành công",
		Data:       user,
	})
}

func (u *UserHandler) SignIn(c echo.Context) error {
	request := req.ReqSignIn{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	user, err := u.UserRepo.CheckSignIn(c.Request().Context(), request)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	// check password
	isTheSame := security.ComparePasswords(user.Password, []byte(request.Password))
	if !isTheSame {
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    "Mật khẩu không đúng",
			Data:       nil,
		})
	}

	// create token
	token, err := security.CreateToken(user.UserID)
	if err != nil {
		return c.JSON(http.StatusForbidden, model.Response{
			StatusCode: http.StatusForbidden,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	saveErr := u.AuthRepo.CreateAuth(user.UserID, token)
	if saveErr != nil {
		return c.JSON(http.StatusForbidden, model.Response{
			StatusCode: http.StatusForbidden,
			Message:    saveErr.Error(),
			Data:       nil,
		})
	}

	user.AccessToken = token.AccessToken
	user.RefreshToken = token.RefreshToken

	// create cookie for client(browser)
	atCookie := &http.Cookie{
		Name:     "access_token",
		Value:    token.AccessToken,
		HttpOnly: true,
		SameSite: 2,
		Expires:  time.Now().Add(time.Minute * 15),
	}
	c.SetCookie(atCookie)

	rtCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    token.RefreshToken,
		SameSite: 2,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 24),
	}
	c.SetCookie(rtCookie)

	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Đăng nhập thành công",
		Data:       user,
	})
}

func (u *UserHandler) Profile(c echo.Context) error {
	tokenAuth, err := security.ExtractAccessTokenMetadata(c.Request())
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	userID, err := u.AuthRepo.FetchAuth(tokenAuth.AccessUUID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	user, err := u.UserRepo.SelectUserByID(c.Request().Context(), userID)
	if err != nil {
		log.Println(err)
		if err == custom_error.UserNotFound {
			return c.JSON(http.StatusNotFound, model.Response{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
				Data:       nil,
			})
		}

		return c.JSON(http.StatusForbidden, model.Response{
			StatusCode: http.StatusForbidden,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Xử lý thành công",
		Data:       user,
	})
}

func (u *UserHandler) UpdateProfile(c echo.Context) error {
	request := req.ReqUpdateUser{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	tokenAuth, err := security.ExtractAccessTokenMetadata(c.Request())
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	userID, err := u.AuthRepo.FetchAuth(tokenAuth.AccessUUID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	hash := security.HashAndSalt([]byte(request.Password))

	user := model.User{
		UserID:   userID,
		Email:    request.Email,
		FullName: request.FullName,
		Password: hash,
	}

	user, err = u.UserRepo.UpdateUser(c.Request().Context(), user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, model.Response{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	return c.JSON(http.StatusCreated, model.Response{
		StatusCode: http.StatusCreated,
		Message:    "Cập nhật thông tin thành công",
		Data:       user,
	})
}

func (u *UserHandler) SignOut(c echo.Context) error {
	extractAt, err := security.ExtractAccessTokenMetadata(c.Request())
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	deleteAtErr := u.AuthRepo.DeleteAccessToken(extractAt.AccessUUID)
	if deleteAtErr != nil {
		log.Println(deleteAtErr)
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    deleteAtErr.Error(),
			Data:       nil,
		})
	}

	extractRt, err := security.ExtractRefreshTokenMetadata(c.Request())
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	deleteRtErr := u.AuthRepo.DeleteRefreshToken(extractRt.RefreshUUID)
	if deleteRtErr != nil {
		log.Println(deleteRtErr)
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    deleteRtErr.Error(),
			Data:       nil,
		})
	}

	atCookie := &http.Cookie{
		Name:   "access_token",
		MaxAge: -1,
	}
	c.SetCookie(atCookie)

	rtCookie := &http.Cookie{
		Name:   "refresh_token",
		MaxAge: -1,
	}
	c.SetCookie(rtCookie)

	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Đăng xuất thành công",
	})
}

func (u *UserHandler) Refresh(c echo.Context) error {
	_, err := c.Cookie("access_token")
	if err != nil {
		extractRt, err := security.ExtractRefreshTokenMetadata(c.Request())
		if err != nil {
			return c.JSON(http.StatusUnauthorized, model.Response{
				StatusCode: http.StatusUnauthorized,
				Message:    "Bạn cần phải đăng nhập",
			})
		}

		deleteErr := u.AuthRepo.DeleteRefreshToken(extractRt.RefreshUUID)
		if deleteErr != nil {
			log.Println(deleteErr)
			return c.JSON(http.StatusUnauthorized, model.Response{
				StatusCode: http.StatusUnauthorized,
				Message:    deleteErr.Error(),
				Data:       nil,
			})
		}

		token, createErr := security.CreateToken(extractRt.UserID)
		if createErr != nil {
			return c.JSON(http.StatusForbidden, model.Response{
				StatusCode: http.StatusForbidden,
				Message:    createErr.Error(),
				Data:       nil,
			})
		}

		saveErr := u.AuthRepo.CreateAuth(extractRt.UserID, token)
		if saveErr != nil {
			return c.JSON(http.StatusForbidden, model.Response{
				StatusCode: http.StatusForbidden,
				Message:    saveErr.Error(),
				Data:       nil,
			})
		}

		tokens := map[string]string{
			"access_token":  token.AccessToken,
			"refresh_token": token.RefreshToken,
		}

		atCookie := &http.Cookie{
			Name:     "access_token",
			Value:    token.AccessToken,
			HttpOnly: true,
			SameSite: 2,
			Expires:  time.Now().Add(time.Minute * 15),
		}
		c.SetCookie(atCookie)

		rtCookie := &http.Cookie{
			Name:     "refresh_token",
			Value:    token.RefreshToken,
			SameSite: 2,
			HttpOnly: true,
			Expires:  time.Now().Add(time.Hour * 24),
		}
		c.SetCookie(rtCookie)

		return c.JSON(http.StatusCreated, model.Response{
			StatusCode: http.StatusCreated,
			Message:    "Xử lý thành công",
			Data:       tokens,
		})
	}

	return c.JSON(http.StatusUnauthorized, model.Response{
		StatusCode: http.StatusUnauthorized,
		Message:    "Access token chưa hết hạn",
	})
}
