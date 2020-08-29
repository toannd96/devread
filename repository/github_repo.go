package repository

import (
	"tech_posts_trending/model"
	"context"
)

type GithubRepo interface {
	SaveRepo(context context.Context, user model.GithubRepo) (model.GithubRepo, error)
	SelectRepos(context context.Context, userId string, limit int) ([]model.GithubRepo, error)
	SelectRepoByName(context context.Context, name string) (model.GithubRepo, error)
	UpdateRepo(context context.Context, user model.GithubRepo) (model.GithubRepo, error)

	//Bookmark
	SelectAllBookmarks(context context.Context, userId string) ([]model.GithubRepo, error)
	Bookmark(context context.Context, bid, nameRepo, userId string) error
	DelBookmark(context context.Context, nameRepo, userId string) error
}
