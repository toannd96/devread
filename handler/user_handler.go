package handler

import (
	"devread/custom_error"
	"devread/helper"
	"devread/model"
	"devread/model/req"
	"devread/repository"
	"devread/security"

	"net/http"
	"net/smtp"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"go.uber.org/zap"
)

type smtpServer struct {
	host string
	port string
}

// Address URI to smtp server.
func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

type UserHandler struct {
	UserRepo repository.UserRepo
	AuthRepo repository.AuthenRepo
	Logger   *zap.Logger
}

// SignUp godoc
// @Summary Create new account
// @Tags user
// @Accept json
// @Produce json
// @Param data body req.ReqSignUp true "user"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 403 {object} model.Response
// @Failure 409 {object} model.Response
// @Router /user/sign-up [post]
func (u *UserHandler) SignUp(c echo.Context) error {
	request := req.ReqSignUp{}
	if err := c.Bind(&request); err != nil {
		u.Logger.Error("Lỗi cú pháp ", zap.Error(err))
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Lỗi cú pháp",
		})
	}

	if err := c.Validate(request); err != nil {
		u.Logger.Error("Lỗi cú pháp ", zap.Error(err))
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Lỗi cú pháp",
		})
	}

	hash := security.HashAndSalt([]byte(request.Password))

	userID, err := uuid.NewUUID()
	if err != nil {
		u.Logger.Error("Tạo mới uuid thất bại ", zap.Error(err))
		return c.JSON(http.StatusForbidden, model.Response{
			StatusCode: http.StatusForbidden,
		})
	}

	user := model.User{
		UserID:   userID.String(),
		FullName: request.FullName,
		Email:    request.Email,
		Password: hash,
		Verify:   false,
	}

	user, err = u.UserRepo.SaveUser(c.Request().Context(), user)
	if err != nil {
		u.Logger.Error("Lưu tài khoản người dùng thất bại ", zap.Error(err))
		return c.JSON(http.StatusConflict, model.Response{
			StatusCode: http.StatusConflict,
		})
	}

	// verify email
	token := helper.CreateTokenHash(user.Email)

	// save token to redis
	saveErr := u.AuthRepo.CreateTokenMail(token, user.UserID)
	if saveErr != nil {
		u.Logger.Error("Tạo token thất bại mail ", zap.Error(saveErr))
		return c.JSON(http.StatusForbidden, model.Response{
			StatusCode: http.StatusForbidden,
		})
	}

	link := "https://devread.herokuapp.com" + "/user/verify?token=" + token

	from := os.Getenv("FROM")
	password := os.Getenv("PASSWORD")
	to := []string{user.Email}

	smtpsv := smtpServer{
		host: os.Getenv("SMTP_HOST"),
		port: os.Getenv("SMTP_PORT"),
	}

	subject := "Xác thực tài khoản\r\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := "Để xác thực tài khoản nhấp vào liên kết <a href='" + link + "'>ở đây</a>."
	message := []byte("Subject:" + subject + mime + "\r\n" + body)

	auth := smtp.PlainAuth("", from, password, smtpsv.host)
	errSendMail := smtp.SendMail(smtpsv.Address(), auth, from, to, message)
	if errSendMail != nil {
		u.Logger.Error("Gửi email thất bại ", zap.Error(errSendMail))
		c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Gửi email thất bại",
		})
	}

	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Tin nhắn xác thực tài khoản được gửi đến email được cung cấp. Vui lòng kiểm tra thư mục thư rác",
	})
}

// ForgotPassword godoc
// @Summary Forgot password
// @Tags user
// @Accept  json
// @Produce  json
// @Param data body req.ReqSignUp true "user"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 401 {object} model.Response
// @Failure 403 {object} model.Response
// @Router /user/password/forgot [post]
func (u *UserHandler) ForgotPassword(c echo.Context) error {
	request := req.ReqSignUp{}
	if err := c.Bind(&request); err != nil {
		u.Logger.Error("Lỗi cú pháp ", zap.Error(err))
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Lỗi cú pháp",
		})
	}

	if err := c.Validate(request); err != nil {
		u.Logger.Error("Lỗi cú pháp ", zap.Error(err))
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Lỗi cú pháp",
		})
	}

	user, err := u.UserRepo.CheckEmail(c.Request().Context(), request)
	if err != nil {
		u.Logger.Error("Email không tồn tại ", zap.Error(err))
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    "Email không tồn tại",
		})
	}

	token := helper.CreateTokenHash(user.Email)

	// save token to redis
	saveErr := u.AuthRepo.CreateTokenMail(token, user.UserID)
	if saveErr != nil {
		u.Logger.Error("Tạo token thất bại mail ", zap.Error(saveErr))
		return c.JSON(http.StatusForbidden, model.Response{
			StatusCode: http.StatusForbidden,
		})
	}

	insertErr := u.AuthRepo.InsertTokenMail(token)
	if insertErr != nil {
		u.Logger.Error("Nhập token mail thất bại ", zap.Error(insertErr))
		return c.JSON(http.StatusForbidden, model.Response{
			StatusCode: http.StatusForbidden,
		})
	}

	link := "https://devread.herokuapp.com" + "/user/password/reset?token=" + token

	from := os.Getenv("FROM")
	password := os.Getenv("PASSWORD")
	to := []string{user.Email}

	smtpsv := smtpServer{
		host: os.Getenv("SMTP_HOST"),
		port: os.Getenv("SMTP_PORT"),
	}

	subject := "Đặt lại mật khẩu\r\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := "Để đặt lại mật khẩu nhấp vào liên kết <a href='" + link + "'>ở đây</a>."
	message := []byte("Subject:" + subject + mime + "\r\n" + body)

	auth := smtp.PlainAuth("", from, password, smtpsv.host)
	errSendMail := smtp.SendMail(smtpsv.Address(), auth, from, to, message)
	if errSendMail != nil {
		u.Logger.Error("Gửi email thất bại ", zap.Error(errSendMail))
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Gửi email thất bại",
		})
	}

	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Tin nhắn được gửi đến email được cung cấp. Vui lòng kiểm tra thư mục thư rác",
	})
}

// VerifyAccount ... godoc
// @Summary Verify email
// @Tags user
// @Accept  json
// @Produce  json
// @Param data body req.PasswordSubmit true "user"
// @Param token query string true "token verify email"
// @Security token-verify-account
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 401 {object} model.Response
// @Failure 403 {object} model.Response
// @Failure 404 {object} model.Response
// @Router /user/verify [post]
func (u *UserHandler) VerifyAccount(c echo.Context) error {
	request := req.PasswordSubmit{}
	if err := c.Bind(&request); err != nil {
		u.Logger.Error("Lỗi cú pháp ", zap.Error(err))
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Lỗi cú pháp",
		})
	}

	if err := c.Validate(request); err != nil {
		u.Logger.Error("Lỗi cú pháp ", zap.Error(err))
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Lỗi cú pháp",
		})
	}

	token := security.ExtractTokenMail(c.Request())

	userID, err := u.AuthRepo.FetchTokenMail(token)
	if err != nil {
		u.Logger.Error("Lỗi khi tìm token mail ", zap.Error(err))
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
		})
	}

	user, err := u.UserRepo.SelectUserByID(c.Request().Context(), userID)
	if err != nil {
		u.Logger.Error("Người dùng không tồn tại ", zap.Error(err))
		return c.JSON(http.StatusNotFound, model.Response{
			StatusCode: http.StatusNotFound,
			Message:    "Người dùng không tồn tại",
		})
	}

	if request.Password != request.Confirm {
		u.Logger.Error("Xác nhận mật khẩu không khớp ", zap.String("Password", request.Password), zap.String("Confirm", request.Confirm))
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    "Xác nhận mật khẩu không khớp",
		})
	}

	// check password
	isTheSame := security.ComparePasswords(user.Password, []byte(request.Password))
	if !isTheSame {
		u.Logger.Error("Mật khẩu không chính xác ", zap.Error(err))
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    "Mật khẩu không đúng",
		})
	}

	user = model.User{
		UserID: userID,
		Verify: true,
	}

	user, err = u.UserRepo.UpdateVerify(c.Request().Context(), user)
	if err != nil {
		u.Logger.Error("Xác thực tài khoản thất bại ", zap.Error(err))
		return c.JSON(http.StatusForbidden, model.Response{
			StatusCode: http.StatusForbidden,
			Message:    "Xác thực tài khoản thất bại",
		})
	}

	deleteAtErr := u.AuthRepo.DeleteTokenMail(token)
	if deleteAtErr != nil {
		u.Logger.Error("Xóa token mail thất bại ", zap.Error(deleteAtErr))
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
		})
	}

	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Xác thực tài khoản thành công",
	})

}

// ResetPassword godoc
// @Summary Reset password
// @Tags user
// @Accept  json
// @Produce  json
// @Param data body req.PasswordSubmit true "user"
// @Param token query string true "token reset password"
// @Security token-reset-password
// @Success 201 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 401 {object} model.Response
// @Failure 403 {object} model.Response
// @Router /user/password/reset [put]
func (u *UserHandler) ResetPassword(c echo.Context) error {
	request := req.PasswordSubmit{}
	if err := c.Bind(&request); err != nil {
		u.Logger.Error("Lỗi cú pháp ", zap.Error(err))
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Lỗi cú pháp",
		})
	}

	if err := c.Validate(request); err != nil {
		u.Logger.Error("Lỗi cú pháp ", zap.Error(err))
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Lỗi cú pháp",
		})
	}

	token := security.ExtractTokenMail(c.Request())

	userID, err := u.AuthRepo.FetchTokenMail(token)
	if err != nil {
		u.Logger.Error("Lỗi khi tìm token mail ", zap.Error(err))
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
		})
	}

	if request.Password != request.Confirm {
		u.Logger.Error("Xác nhận mật khẩu không khớp ", zap.String("Password", request.Password), zap.String("Confirm", request.Confirm))
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    "Xác nhận mật khẩu không khớp",
		})
	}

	hash := security.HashAndSalt([]byte(request.Password))

	user := model.User{
		UserID:   userID,
		Password: hash,
	}

	user, err = u.UserRepo.UpdatePassword(c.Request().Context(), user)
	if err != nil {
		u.Logger.Error("Cập nhật mật khẩu thất bại ", zap.Error(err))
		return c.JSON(http.StatusForbidden, model.Response{
			StatusCode: http.StatusForbidden,
			Message:    "Cập nhật mật khẩu thất bại",
		})
	}

	deleteAtErr := u.AuthRepo.DeleteTokenMail(token)
	if deleteAtErr != nil {
		u.Logger.Error("Xóa token mail thất bại ", zap.Error(deleteAtErr))
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
		})
	}

	return c.JSON(http.StatusCreated, model.Response{
		StatusCode: http.StatusCreated,
		Message:    "Tạo mới mật khẩu thành công",
	})
}

// SignIn godoc
// @Summary Sign in to access your account
// @Tags user
// @Accept  json
// @Produce  json
// @Param data body req.ReqSignIn true "user"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 401 {object} model.Response
// @Failure 403 {object} model.Response
// @Failure 404 {object} model.Response
// @Router /user/sign-in [post]
func (u *UserHandler) SignIn(c echo.Context) error {
	request := req.ReqSignIn{}
	if err := c.Bind(&request); err != nil {
		u.Logger.Error("Lỗi cú pháp ", zap.Error(err))
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Lỗi cú pháp",
		})
	}

	if err := c.Validate(request); err != nil {
		u.Logger.Error("Lỗi cú pháp ", zap.Error(err))
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Lỗi cú pháp",
		})
	}

	user, err := u.UserRepo.CheckSignIn(c.Request().Context(), request)
	if err != nil {
		u.Logger.Error("Tài khoản không tồn tại ", zap.Error(err))
		return c.JSON(http.StatusNotFound, model.Response{
			StatusCode: http.StatusNotFound,
			Message:    "Tài khoản không tồn tại",
		})
	}

	if !user.Verify {
		u.Logger.Debug("Tài khoản chưa được xác thực")
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    "Tài khoản chưa được xác thực",
		})
	}

	// check password
	isTheSame := security.ComparePasswords(user.Password, []byte(request.Password))
	if !isTheSame {
		u.Logger.Error("Mật khẩu không chính xác ", zap.Error(err))
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    "Mật khẩu không chính xác",
		})
	}

	// create token
	token, err := security.CreateToken(user)
	if err != nil {
		u.Logger.Error("Tạo token thất bại ", zap.Error(err))
		return c.JSON(http.StatusForbidden, model.Response{
			StatusCode: http.StatusForbidden,
		})
	}
	user.Token = token

	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Đăng nhập thành công",
		Data:       user,
	})
}

// Profile godoc
// @Summary Get user profile
// @Tags profile
// @Accept  json
// @Produce  json
// @Security jwt
// @Success 200 {object} model.Response
// @Failure 403 {object} model.Response
// @Failure 404 {object} model.Response
// @Router /user/profile [get]
func (u *UserHandler) Profile(c echo.Context) error {
	tokenData := c.Get("user").(*jwt.Token)
	claims := tokenData.Claims.(*model.TokenDetails)

	user, err := u.UserRepo.SelectUserByID(c.Request().Context(), claims.UserID)
	if err != nil {
		u.Logger.Error("Người dùng không tồn tại ", zap.Error(err))
		if err == custom_error.UserNotFound {
			return c.JSON(http.StatusNotFound, model.Response{
				StatusCode: http.StatusNotFound,
				Message:    "Người dùng không tồn tại",
			})
		}

		return c.JSON(http.StatusForbidden, model.Response{
			StatusCode: http.StatusForbidden,
			Message:    "Truy cập không được phép",
		})
	}

	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Xử lý thành công",
		Data:       user,
	})
}

// UpdateProfile godoc
// @Summary Update user profile
// @Tags profile
// @Accept  json
// @Produce  json
// @Security jwt
// @Param data body req.ReqUpdateUser true "user"
// @Success 201 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 401 {object} model.Response
// @Failure 422 {object} model.Response
// @Router /user/profile/update [put]
func (u *UserHandler) UpdateProfile(c echo.Context) error {
	request := req.ReqUpdateUser{}
	if err := c.Bind(&request); err != nil {
		u.Logger.Error("Lỗi cú pháp ", zap.Error(err))
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Lỗi cú pháp",
		})
	}

	if err := c.Validate(request); err != nil {
		u.Logger.Error("Lỗi cú pháp ", zap.Error(err))
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Lỗi cú pháp",
		})
	}

	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*model.TokenDetails)

	if request.Password != request.Confirm {
		u.Logger.Error("Xác nhận mật khẩu không khớp ", zap.String("Password", request.Password), zap.String("Confirm", request.Confirm))
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    "Xác nhận mật khẩu không khớp",
		})
	}

	hash := security.HashAndSalt([]byte(request.Password))

	var err error
	if request.FullName == "" {
		if len(request.Password) < 8 {
			return c.JSON(http.StatusBadRequest, model.Response{
				StatusCode: http.StatusBadRequest,
				Message:    "Mật khẩu tối thiểu 8 ký tự",
			})
		}
		user := model.User{
			UserID:   claims.UserID,
			Password: hash,
		}

		user, err = u.UserRepo.UpdateUser(c.Request().Context(), user)
		if err != nil {
			u.Logger.Error("Cập nhật thông tin người dùng thất bại ", zap.Error(err))
			return c.JSON(http.StatusUnprocessableEntity, model.Response{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Cập nhật thông tin người dùng thất bại",
			})
		}
	}

	if request.Password == "" {
		user := model.User{
			UserID:   claims.UserID,
			FullName: request.FullName,
		}

		user, err = u.UserRepo.UpdateUser(c.Request().Context(), user)
		if err != nil {
			u.Logger.Error("Cập nhật thông tin người dùng thất bại ", zap.Error(err))
			return c.JSON(http.StatusUnprocessableEntity, model.Response{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    "Cập nhật thông tin người dùng thất bại",
			})
		}
	}

	user := model.User{
		UserID:   claims.UserID,
		FullName: request.FullName,
		Password: hash,
	}

	user, err = u.UserRepo.UpdateUser(c.Request().Context(), user)
	if err != nil {
		u.Logger.Error("Cập nhật thông tin người dùng thất bại ", zap.Error(err))
		return c.JSON(http.StatusUnprocessableEntity, model.Response{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Cập nhật thông tin người dùng thất bại",
		})
	}

	return c.JSON(http.StatusCreated, model.Response{
		StatusCode: http.StatusCreated,
		Message:    "Cập nhật thông tin thành công",
	})
}
