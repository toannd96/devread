package repo_impl

import (
	"context"
	"database/sql"
	"github.com/labstack/gommon/log"
	"github.com/lib/pq"
	"tech_posts_trending/custom_error"
	"tech_posts_trending/db"
	"tech_posts_trending/model"
	"tech_posts_trending/repository"
	"time"
)

type PostRepoImpl struct {
	sql *db.Sql
}

func NewPostRepo(sql *db.Sql) repository.PostRepo {
	return &PostRepoImpl{
		sql: sql,
	}
}

func (p PostRepoImpl) SavePost(context context.Context, post model.Post) (model.Post, error) {
	statement := `INSERT INTO posts(name, link, tags) 
          		  VALUES(:name,:link, :tags)`
	_, err := p.sql.Db.NamedExecContext(context, statement, post)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" {
				return post, custom_error.PostConflict
			}
		}
		log.Error(err.Error())
		return post, custom_error.PostInsertFail
	}

	return post, nil
}

func (p PostRepoImpl) SelectPostByName(context context.Context, name string) (model.Post, error) {
	var post = model.Post{}
	err := p.sql.Db.GetContext(context, &post,
		`SELECT * FROM posts WHERE name=$1`, name)

	if err != nil {
		if err == sql.ErrNoRows {
			return post, custom_error.PostNotFound
		}
		log.Error(err.Error())
		return post, err
	}
	return post, nil
}

func (p PostRepoImpl) SelectPostByTags(context context.Context, tags string) ([]model.Post, error) {
	var posts = []model.Post{}
	err := p.sql.Db.SelectContext(context, &posts,
		`SELECT * FROM posts WHERE tags in ($1)`, tags)

	if err != nil {
		if err == sql.ErrNoRows {
			return posts, custom_error.PostNotFound
		}
		log.Error(err.Error())
		return posts, err
	}
	return posts, nil
}

func (p PostRepoImpl) UpdatePost(context context.Context, post model.Post) (model.Post, error) {
	sqlStatement := `
		UPDATE posts
		SET
			link = :link,
			tags = :tags
		WHERE name = :name
	`

	result, err := p.sql.Db.NamedExecContext(context, sqlStatement, post)
	if err != nil {
		log.Error(err.Error())
		return post, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		log.Error(err.Error())
		return post, custom_error.PostNotUpdated
	}
	if count == 0 {
		return post, custom_error.PostNotUpdated
	}

	return post, nil
}

func (p PostRepoImpl) SelectAllPost(context context.Context) ([]model.Post, error) {
	posts := []model.Post{}
	err := p.sql.Db.SelectContext(context, &posts,
		`SELECT * FROM posts`)

	if err != nil {
		if err == sql.ErrNoRows {
			return posts, custom_error.PostNotFound
		}
		log.Error(err.Error())
		return posts, err
	}
	return posts, nil
}

func (p PostRepoImpl) SelectAllBookmark(context context.Context, userId string) ([]model.Post, error) {
	posts := []model.Post{}
	err := p.sql.Db.SelectContext(context, &posts,
		`SELECT 
					posts.name, posts.link, posts.tags
				FROM bookmarks 
				INNER JOIN posts
				ON bookmarks.user_id=$1 AND posts.name = bookmarks.post_name`, userId)

	if err != nil {
		if err == sql.ErrNoRows {
			return posts, custom_error.BookmarkNotFound
		}
		log.Error(err.Error())
		return posts, err
	}
	return posts, nil
}

func (p PostRepoImpl) Bookmark(context context.Context, bookmarkId, namePost, userId string) error {
	statement := `INSERT INTO bookmarks(
					bookmark_id, user_id, post_name, created_at, updated_at) 
          		  VALUES($1, $2, $3, $4, $5)`

	now := time.Now()
	_, err := p.sql.Db.ExecContext(
		context, statement, bookmarkId, userId,
		namePost, now, now)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" {
				return custom_error.BookmarkConflic
			}
		}
		log.Error(err.Error())
		return custom_error.BookmarkFail
	}

	return nil
}

func (p PostRepoImpl) DelBookmark(context context.Context, namePost, userId string) error {
	result := p.sql.Db.MustExecContext(
		context,
		"DELETE FROM bookmarks WHERE post_name = $1 AND user_id = $2",
		namePost, userId)

	index, err := result.RowsAffected()
	if index == 0 {
		return custom_error.BookmarkNotFound
	}
	if err != nil {
		log.Error(err.Error())
		return custom_error.DelBookmarkFail
	}

	return nil
}
