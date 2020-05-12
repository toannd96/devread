package custom_error

import "errors"

var (
	//repo
	RepoNotUpdated = errors.New("Cập nhật thông tin Repo thất bại")
	RepoNotFound   = errors.New("Repo không tồn tại")
	RepoConflict   = errors.New("Repo đã tồn tại")
	RepoInsertFail = errors.New("Thêm Repo thất bại")

	//bookmark
	BookmarkNotFound = errors.New("Bookmark không tồn tại")
	BookmarkFail     = errors.New("Bookmark thất bại")
	DelBookmarkFail  = errors.New("DelBookmark thất bại")
	BookmarkConflic  = errors.New("Bookmark đã tồn tại")

	//genneral
	ErrorSql = errors.New("Lỗi SQL")
)
