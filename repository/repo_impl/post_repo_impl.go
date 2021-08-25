package repo_impl

import (
	"context"
	"database/sql"

	"devread/custom_error"
	"devread/db"
	"devread/model"
	"devread/repository"

	"github.com/lib/pq"
)

type PostRepoImpl struct {
	sql *db.Sql
}

func NewPostRepo(sql *db.Sql) repository.PostRepo {
	return &PostRepoImpl{
		sql: sql,
	}
}

func (p PostRepoImpl) Save(context context.Context, post model.Post) (model.Post, error) {
	statement := `INSERT INTO posts(name, link, tag) 
          		  VALUES(:name, :link, :tag)`
	_, err := p.sql.Db.NamedExecContext(context, statement, post)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" {
				return post, custom_error.PostConflict
			}
		}
		return post, custom_error.PostInsertFail
	}
	return post, nil
}

func (p PostRepoImpl) SelectByLink(context context.Context, link string) (model.Post, error) {
	var post = model.Post{}
	err := p.sql.Db.GetContext(context, &post,
		`SELECT * FROM posts WHERE link=$1`, link)
	if err != nil {
		if err == sql.ErrNoRows {
			return post, custom_error.PostNotFound
		}
		return post, err
	}
	return post, nil
}

func (p PostRepoImpl) SelectByTag(context context.Context, tag string) ([]model.Post, error) {
	var posts = []model.Post{}
	err := p.sql.Db.SelectContext(context, &posts,
		`SELECT * FROM posts WHERE tag=$1`, tag)

	if err != nil {
		if err == sql.ErrNoRows {
			return posts, custom_error.PostNotFound
		}
		return posts, err
	}
	return posts, nil
}

func (p PostRepoImpl) Update(context context.Context, post model.Post) (model.Post, error) {
	sqlStatement := `
		UPDATE posts
		SET
		    name = :name,
		WHERE link = :link
	`
	result, err := p.sql.Db.NamedExecContext(context, sqlStatement, post)
	if err != nil {
		return post, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return post, custom_error.PostNotUpdated
	}
	if count == 0 {
		return post, custom_error.PostNotUpdated
	}
	return post, nil
}

func (p PostRepoImpl) SelectAll(context context.Context) ([]model.Post, error) {
	posts := []model.Post{}
	err := p.sql.Db.SelectContext(context, &posts,
		`SELECT * FROM posts`)
	if err != nil {
		if err == sql.ErrNoRows {
			return posts, custom_error.PostNotFound
		}
		return posts, err
	}
	return posts, nil
}
