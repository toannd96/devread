package handler

import (
	"tech_posts_trending/model"
	"tech_posts_trending/model/req"
	"tech_posts_trending/log"
	"tech_posts_trending/repository"
	"tech_posts_trending/security"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RepoHandler struct {
	GithubRepo repository.GithubRepo
	AuthRepo   repository.AuthRepo
}

// RepoTrending godoc
// @Summary Get repository github trending by user
// @Tags repository-service
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Response
// @Failure 401 {object} model.Response
// @Router /user/github/trending [get]
func (r *RepoHandler) RepoTrending(c echo.Context) error {
	tokenAuth, err := security.ExtractAccessTokenMetadata(c.Request())
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    err.Error(),
		})
	}

	userID, err := r.AuthRepo.FetchAuth(tokenAuth.AccessUUID)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    "Truy cập không được phép",
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

// SelectBookmarks godoc
// @Summary Get list bookmark
// @Tags bookmark-service
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Response
// @Failure 401 {object} model.Response
// @Router /user/bookmark/list [get]
func (r *RepoHandler) SelectBookmarks(c echo.Context) error {
	tokenAuth, err := security.ExtractAccessTokenMetadata(c.Request())
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    err.Error(),
		})
	}

	userID, err := r.AuthRepo.FetchAuth(tokenAuth.AccessUUID)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    "Truy cập không được phép",
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

// Bookmark godoc
// @Summary Add bookmark
// @Tags bookmark-service
// @Accept  json
// @Produce  json
// @Param data body req.ReqBookmark true "user"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 401 {object} model.Response
// @Failure 403 {object} model.Response
// @Failure 409 {object} model.Response
// @Router /user/bookmark/add [post]
func (r *RepoHandler) Bookmark(c echo.Context) error {
	req := req.ReqBookmark{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Lỗi cú pháp",
		})
	}

	// validate thông tin gửi lên
	err := c.Validate(req)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Lỗi cú pháp",
		})
	}

	tokenAuth, err := security.ExtractAccessTokenMetadata(c.Request())
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    err.Error(),
		})
	}

	userID, err := r.AuthRepo.FetchAuth(tokenAuth.AccessUUID)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    "Truy cập không được phép",
		})
	}

	bId, err := uuid.NewUUID()
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusForbidden, model.Response{
			StatusCode: http.StatusForbidden,
			Message:    err.Error(),
		})
	}

	err = r.GithubRepo.Bookmark(
		c.Request().Context(),
		bId.String(),
		req.RepoName,
		userID)

	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusConflict, model.Response{
			StatusCode: http.StatusConflict,
			Message:    "Bookmark thất bại",
		})
	}

	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Bookmark thành công",
	})
}

// DelBookmark godoc
// @Summary Delete bookmark
// @Tags bookmark-service
// @Accept  json
// @Produce  json
// @Param data body req.ReqBookmark true "user"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 401 {object} model.Response
// @Failure 409 {object} model.Response
// @Router /user/bookmark/delete [delete]
func (r *RepoHandler) DelBookmark(c echo.Context) error {
	req := req.ReqBookmark{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Lỗi cú pháp",
		})
	}

	// validate thông tin gửi lên
	err := c.Validate(req)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Lỗi cú pháp",
		})
	}

	tokenAuth, err := security.ExtractAccessTokenMetadata(c.Request())
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    err.Error(),
		})
	}

	userID, err := r.AuthRepo.FetchAuth(tokenAuth.AccessUUID)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusUnauthorized, model.Response{
			StatusCode: http.StatusUnauthorized,
			Message:    "Truy cập không được phép",
		})
	}

	err = r.GithubRepo.DelBookmark(
		c.Request().Context(),
		req.RepoName, userID)

	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusConflict, model.Response{
			StatusCode: http.StatusConflict,
			Message:    "Bookmark không tồn tại",
		})
	}

	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Xoá bookmark thành công",
	})
}
