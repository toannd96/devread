package repository

import (
	"context"
	"devread/model"
)

type BookmarkRepo interface {
	SelectAll(context context.Context, userId string) ([]model.Post, error)
	Bookmark(context context.Context, bid, namePost, userId string) error
	Delete(context context.Context, namePost, userId string) error
}
