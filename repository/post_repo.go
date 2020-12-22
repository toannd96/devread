package repository

import (
	"context"
	"tech_posts_trending/model"
)

type PostRepo interface {
	// Post
	UpdatePost(context context.Context, post model.Post) (model.Post, error)
	SavePost(context context.Context, post model.Post) (model.Post, error)
	SelectAllPost(context context.Context) ([]model.Post, error)
	SelectPostByName(context context.Context, name string) (model.Post, error)
	SelectPostByTag(context context.Context, tag string) ([]model.Post, error)

	// Bookmark
	SelectAllBookmark(context context.Context, userId string) ([]model.Post, error)
	Bookmark(context context.Context, bid, namePost, userId string) error
	DelBookmark(context context.Context, namePost, userId string) error
}
