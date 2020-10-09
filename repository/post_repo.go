package repository

import (
	"tech_posts_trending/model"
	"context"
)

type PostRepo interface {
	// Post
	SavePost(context context.Context, post model.Post) (model.Post, error)
	SelectAllPost(context context.Context) ([]model.Post, error)
	SelectPostByName(context context.Context, name string) (model.Post, error)
	UpdatePost(context context.Context, post model.Post) (model.Post, error)

	// Bookmark
	SelectAllBookmark(context context.Context, userId string) ([]model.Post, error)
	Bookmark(context context.Context, bid, namePost, userId string) error
	DelBookmark(context context.Context, namePost, userId string) error
}
