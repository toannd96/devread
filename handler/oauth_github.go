package handler

import (
	"backend-viblo-trending/model"
	"backend-viblo-trending/repository"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dchest/uniuri"
	"github.com/labstack/echo"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

type OauthGithub struct {
	UserRepo repository.UserRepo
	AuthRepo repository.AuthRepo
}

const (
	apiGithubProfileURL = "https://api.github.com/user"
	redirectURL         = "http://localhost:3000/user/github/callback"
)

// Trao đổi mã ủy quyền cho mã thông báo truy cập
var (
	githubOauthConfig = &oauth2.Config{
		ClientID:     "57a3ce2aa87dc8425485",
		ClientSecret: "fdd78fe0c1da848fa27b62e7b38a14e8043b8fc5",
		RedirectURL:  redirectURL,
		Scopes:       []string{"user:email"},
		Endpoint:     githuboauth.Endpoint,
	}
	// Chuỗi ngẫu nhiên cho các lệnh gọi API oauth2 để bảo vệ chống lại CSRF
	oauthGitStateString = uniuri.New()
)

func (oau *OauthGithub) GithubSignIn(c echo.Context) error {
	url := githubOauthConfig.AuthCodeURL(oauthGitStateString)
	fmt.Println(url)
	return c.Redirect(307, url)
}

func (oau *OauthGithub) GithubCallback(c echo.Context) error {

	state := c.FormValue("state")
	if state != oauthGitStateString {
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    "Mã thông báo không hợp lệ",
			Data:       nil,
		})
	}

	code := c.FormValue("code")
	token, err := githubOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Trao đổi mã thông báo thất bại",
			Data:       nil,
		})
	}

	// Lấy thông tin hồ sơ về người dùng hiện tại.
	httpClient := &http.Client{}
	req, _ := http.NewRequest("GET", apiGithubProfileURL, nil)
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	res, err := httpClient.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "Không tìm thấy hồ sơ người dùng",
			Data:       nil,
		})
	}
	defer res.Body.Close()

	user := make(map[string]string)
	json.NewDecoder(res.Body).Decode(&user)

	userInfo := &model.GithubProfile{
		Email:    user["email"],
		FullName: user["name"],
	}

	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Xử lý thành công",
		Data:       userInfo,
	})
}
