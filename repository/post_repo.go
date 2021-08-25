package repository

import (
	"context"

	"devread/model"
)

type PostRepo interface {
	Update(context context.Context, post model.Post) (model.Post, error)
	Save(context context.Context, post model.Post) (model.Post, error)
	SelectAll(context context.Context) ([]model.Post, error)
	SelectByTag(context context.Context, tag string) ([]model.Post, error)
	SelectByLink(context context.Context, link string) (model.Post, error)
}
