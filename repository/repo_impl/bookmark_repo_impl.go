package repo_impl

import (
	"context"
	"database/sql"
	"time"

	"devread/custom_error"
	"devread/db"
	"devread/model"
	"devread/repository"

	"github.com/lib/pq"
)

type BookmarkRepoImpl struct {
	sql *db.Sql
}

func NewBookmarkRepo(sql *db.Sql) repository.BookmarkRepo {
	return &BookmarkRepoImpl{
		sql: sql,
	}
}

func (b BookmarkRepoImpl) SelectAll(context context.Context, userId string) ([]model.Post, error) {
	posts := []model.Post{}
	err := b.sql.Db.SelectContext(context, &posts,
		`SELECT 
					posts.name, posts.link, posts.tag
				FROM bookmarks 
				INNER JOIN posts
				ON bookmarks.user_id=$1 AND posts.name = bookmarks.post_name`, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return posts, custom_error.BookmarkNotFound
		}
		return posts, err
	}
	return posts, nil
}

func (b BookmarkRepoImpl) Bookmark(context context.Context, bookmarkId, namePost, userId string) error {
	statement := `INSERT INTO bookmarks(
					bookmark_id, user_id, post_name, created_at, updated_at) 
          		  VALUES($1, $2, $3, $4, $5)`
	now := time.Now()
	_, err := b.sql.Db.ExecContext(
		context, statement, bookmarkId, userId,
		namePost, now, now)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" {
				return custom_error.BookmarkConflic
			}
		}
		return custom_error.BookmarkFail
	}
	return nil
}

func (b BookmarkRepoImpl) Delete(context context.Context, namePost, userId string) error {
	result := b.sql.Db.MustExecContext(
		context,
		"DELETE FROM bookmarks WHERE post_name = $1 AND user_id = $2",
		namePost, userId)

	index, err := result.RowsAffected()
	if index == 0 {
		return custom_error.BookmarkNotFound
	}
	if err != nil {
		return custom_error.DelBookmarkFail
	}
	return nil
}
