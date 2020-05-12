package handler

import (
	"backend-viblo-trending/model"
	"backend-viblo-trending/model/req"
	"backend-viblo-trending/repository"
	"backend-viblo-trending/security"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"strings"
)

type RepoHandler struct {
	GithubRepo repository.GithubRepo
	AuthRepo repository.AuthRepo
}

func (r *RepoHandler) RepoTrending(c echo.Context) error {
	tokenAuth, err := security.ExtractAccessTokenMetadata(c.Request())
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	userID, err := r.AuthRepo.FetchAuth(tokenAuth.AccessUUID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    "Truy cập không được phép",
			Data:       nil,
		})
	}

	repos, _ := r.GithubRepo.SelectRepos(c.Request().Context(), userID, 25)
	for i, repo := range repos {
		repos[i].Contributors = strings.Split(repo.BuildBy, ",")
	}

	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Xử lý thành công",
		Data:       repos,
	})
}

func (r *RepoHandler) SelectBookmarks(c echo.Context) error {
	tokenAuth, err := security.ExtractAccessTokenMetadata(c.Request())
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	userID, err := r.AuthRepo.FetchAuth(tokenAuth.AccessUUID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    "Truy cập không được phép",
			Data:       nil,
		})
	}

	repos, _ := r.GithubRepo.SelectAllBookmarks(
		c.Request().Context(),
		userID)

	for i, repo := range repos {
		repos[i].Contributors = strings.Split(repo.BuildBy, ",")
	}

	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Xử lý thành công",
		Data:       repos,
	})
}

func (r *RepoHandler) Bookmark(c echo.Context) error {
	req := req.ReqBookmark{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	// validate thông tin gửi lên
	err := c.Validate(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
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

	userID, err := r.AuthRepo.FetchAuth(tokenAuth.AccessUUID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    "Truy cập không được phép",
			Data:       nil,
		})
	}

	bId, err := uuid.NewUUID()
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusForbidden, model.Response{
			StatusCode: http.StatusForbidden,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	err = r.GithubRepo.Bookmark(
		c.Request().Context(),
		bId.String(),
		req.RepoName,
		userID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Bookmark thành công",
		Data:       nil,
	})
}

func (r *RepoHandler) DelBookmark(c echo.Context) error {
	req := req.ReqBookmark{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	// validate thông tin gửi lên
	err := c.Validate(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
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

	userID, err := r.AuthRepo.FetchAuth(tokenAuth.AccessUUID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    "Truy cập không được phép",
			Data:       nil,
		})
	}

	err = r.GithubRepo.DelBookmark(
		c.Request().Context(),
		req.RepoName, userID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Xoá bookmark thành công",
		Data:       nil,
	})
}
