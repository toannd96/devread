package handler

import (
	"devread/model"
	"devread/model/req"
	"devread/repository"

	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"go.uber.org/zap"
)

func GetQueryTag(r *http.Request) string {
	tag := r.URL.Query().Get("tag")
	return tag
}

type PostHandler struct {
	PostRepo     repository.PostRepo
	AuthRepo     repository.AuthenRepo
	BookmarkRepo repository.BookmarkRepo
	Logger       *zap.Logger
}

// PostTrending godoc
// @Summary Get all posts trending
// @Tags post-service
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Response
// @Failure 401 {object} model.Response
// @Router /trend [get]
func (post *PostHandler) PostTrending(c echo.Context) error {
	repos, err := post.PostRepo.SelectAll(c.Request().Context())
	if err != nil {
		post.Logger.Error("Lỗi khi chọn tất cả bài đăng thịnh hành ", zap.Error(err))
		return c.JSON(http.StatusNotFound, model.Response{
			StatusCode: http.StatusNotFound,
			Message:    "Nhận tất cả các bài đăng thịnh hành thất bại",
		})
	}
	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Xử lý thành công",
		Data:       repos,
	})
}

// SearchPost godoc
// @Summary Search post by tag
// @Tags post-service
// @Accept  json
// @Produce  json
// @Param tag query string true "tag of posts"
// @Success 200 {object} model.Response
// @Failure 404 {object} model.Response
// @Router /posts [get]
func (post *PostHandler) SearchPost(c echo.Context) error {
	repos, err := post.PostRepo.SelectByTag(c.Request().Context(), GetQueryTag(c.Request()))
	if err != nil {
		post.Logger.Error("Không tìm thấy bài viết theo tag ", zap.Error(err))
		return c.JSON(http.StatusNotFound, model.Response{
			StatusCode: http.StatusNotFound,
			Message:    "Không tìm thấy bài viết",
		})
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
// @Security jwt
// @Success 200 {object} model.Response
// @Failure 401 {object} model.Response
// @Router /user/bookmark/list [get]
func (post *PostHandler) SelectBookmarks(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*model.TokenDetails)

	repos, err := post.BookmarkRepo.SelectAll(c.Request().Context(), claims.UserID)
	if err != nil {
		post.Logger.Error("Lỗi khi chọn tất cả dấu trang ", zap.Error(err))
		return c.JSON(http.StatusNotFound, model.Response{
			StatusCode: http.StatusNotFound,
			Message:    "Không tìm thấy bài viết",
		})
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
// @Security jwt
// @Param data body req.ReqBookmark true "user"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 401 {object} model.Response
// @Failure 403 {object} model.Response
// @Failure 409 {object} model.Response
// @Router /user/bookmark/add [post]
func (post *PostHandler) Bookmark(c echo.Context) error {
	req := req.ReqBookmark{}
	if err := c.Bind(&req); err != nil {
		post.Logger.Error("Lỗi cú pháp ", zap.Error(err))
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Lỗi cú pháp",
		})
	}

	// validate thông tin gửi lên
	err := c.Validate(req)
	if err != nil {
		post.Logger.Error("Lỗi cú pháp ", zap.Error(err))
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Lỗi cú pháp",
		})
	}
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*model.TokenDetails)

	bId, err := uuid.NewUUID()
	if err != nil {
		post.Logger.Error("Tạo mới uuid thất bại ", zap.Error(err))
		return c.JSON(http.StatusForbidden, model.Response{
			StatusCode: http.StatusForbidden,
		})
	}

	err = post.BookmarkRepo.Bookmark(
		c.Request().Context(),
		bId.String(),
		req.PostName,
		claims.UserID)

	if err != nil {
		post.Logger.Error("Đánh dấu repo mới thất bại ", zap.Error(err))
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
// @Security jwt
// @Param data body req.ReqBookmark true "user"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 401 {object} model.Response
// @Failure 409 {object} model.Response
// @Router /user/bookmark/delete [delete]
func (post *PostHandler) DelBookmark(c echo.Context) error {
	req := req.ReqBookmark{}
	if err := c.Bind(&req); err != nil {
		post.Logger.Error("Lỗi cú pháp ", zap.Error(err))
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Lỗi cú pháp",
		})
	}

	// validate thông tin gửi lên
	err := c.Validate(req)
	if err != nil {
		post.Logger.Error("Lỗi cú pháp ", zap.Error(err))
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    "Lỗi cú pháp",
		})
	}

	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*model.TokenDetails)

	err = post.BookmarkRepo.Delete(
		c.Request().Context(),
		req.PostName, claims.UserID)

	if err != nil {
		post.Logger.Error("Lỗi khi xóa bookmark ", zap.Error(err))
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
