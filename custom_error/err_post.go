package custom_error

import "errors"

var (
	//post
	PostNotUpdated = errors.New("Cập nhật thông tin bài viết thất bại")
	PostNotFound   = errors.New("Bài viết không tồn tại")
	PostConflict   = errors.New("Bài viết đã tồn tại")
	PostInsertFail = errors.New("Thêm bài viết thất bại")

	//bookmark
	BookmarkNotFound = errors.New("Bookmark không tồn tại")
	BookmarkFail     = errors.New("Bookmark thất bại")
	DelBookmarkFail  = errors.New("DelBookmark thất bại")
	BookmarkConflic  = errors.New("Bookmark đã tồn tại")

	//genneral
	ErrorSql = errors.New("Lỗi SQL")
)
