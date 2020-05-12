package repo_impl

import (
	"backend-viblo-trending/custom_error"
	"backend-viblo-trending/db"
	"backend-viblo-trending/model"
	"backend-viblo-trending/repository"
	"context"
	"database/sql"
	"github.com/labstack/gommon/log"
	"github.com/lib/pq"
	"time"
)

type GithubRepoImpl struct {
	sql *db.Sql
}

func NewGithubRepo(sql *db.Sql) repository.GithubRepo {
	return &GithubRepoImpl{
		sql: sql,
	}
}

func (g GithubRepoImpl) SelectRepoByName(context context.Context, name string) (model.GithubRepo, error) {
	var repo = model.GithubRepo{}
	err := g.sql.Db.GetContext(context, &repo,
		`SELECT * FROM repos WHERE name=$1`, name)

	if err != nil {
		if err == sql.ErrNoRows {
			return repo, custom_error.RepoNotFound
		}
		log.Error(err.Error())
		return repo, err
	}
	return repo, nil
}

func (g GithubRepoImpl) SaveRepo(context context.Context, repo model.GithubRepo) (model.GithubRepo, error) {
	// name, description, url, color, lang, fork, stars, stars_today, build_by, created_at, updated_at
	statement := `INSERT INTO repos(
					name, description, url, color, lang, fork, stars, 
 			        stars_today, build_by, created_at, updated_at) 
          		  VALUES(
					:name,:description, :url, :color, :lang, :fork, :stars, 
					:stars_today, :build_by, :created_at, :updated_at
				  )`

	repo.CreatedAt = time.Now()
	repo.UpdatedAt = time.Now()

	_, err := g.sql.Db.NamedExecContext(context, statement, repo)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" {
				return repo, custom_error.RepoConflict
			}
		}
		log.Error(err.Error())
		return repo, custom_error.RepoInsertFail
	}

	return repo, nil
}

func (g GithubRepoImpl) SelectRepos(context context.Context, userId string, limit int) ([]model.GithubRepo, error) {
	var repos []model.GithubRepo
	err := g.sql.Db.SelectContext(context, &repos,
		`
			SELECT 
				repos.name, repos.description, repos.url, repos.color, repos.lang, 
				repos.fork, repos.stars, repos.stars_today, repos.build_by, repos.updated_at, 
				COALESCE(repos.name = bookmarks.repo_name, FALSE) as bookmarked
			FROM repos
			FULL OUTER JOIN bookmarks 
			ON repos.name = bookmarks.repo_name AND 
			   bookmarks.user_id=$1  
			WHERE repos.name IS NOT NULL 
			ORDER BY updated_at ASC LIMIT $2
		`, userId, limit)

	if err != nil {
		log.Error(err.Error())
		return repos, err
	}

	return repos, nil
}

func (g GithubRepoImpl) UpdateRepo(context context.Context, repo model.GithubRepo) (model.GithubRepo, error) {
	// name, description, url, color, lang, fork, stars, stars_today, build_by, created_at, updated_at
	sqlStatement := `
		UPDATE repos
		SET 
			stars  = :stars,
			fork = :fork,
			stars_today = :stars_today,
			build_by = :build_by,
			updated_at = :updated_at
		WHERE name = :name
	`

	result, err := g.sql.Db.NamedExecContext(context, sqlStatement, repo)
	if err != nil {
		log.Error(err.Error())
		return repo, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		log.Error(err.Error())
		return repo, custom_error.RepoNotUpdated
	}
	if count == 0 {
		return repo, custom_error.RepoNotUpdated
	}

	return repo, nil
}

func (g GithubRepoImpl) SelectAllBookmarks(context context.Context, userId string) ([]model.GithubRepo, error) {
	repos := []model.GithubRepo{}
	err := g.sql.Db.SelectContext(context, &repos,
		`SELECT 
					repos.name, repos.description, repos.url, 
					repos.color, repos.lang, repos.fork, repos.stars, 
					repos.stars_today, repos.build_by, true as bookmarked
				FROM bookmarks 
				INNER JOIN repos
				ON bookmarks.user_id=$1 AND repos.name = bookmarks.repo_name`, userId)

	if err != nil {
		if err == sql.ErrNoRows {
			return repos, custom_error.BookmarkNotFound
		}
		log.Error(err.Error())
		return repos, err
	}
	return repos, nil
}

func (g GithubRepoImpl) Bookmark(context context.Context, bid, nameRepo, userId string) error {
	statement := `INSERT INTO bookmarks(
					bid, user_id, repo_name, created_at, updated_at) 
          		  VALUES($1, $2, $3, $4, $5)`

	now := time.Now()
	_, err := g.sql.Db.ExecContext(
		context, statement, bid, userId,
		nameRepo, now, now)

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

func (g GithubRepoImpl) DelBookmark(context context.Context, nameRepo, userId string) error {
	result := g.sql.Db.MustExecContext(
		context,
		"DELETE FROM bookmarks WHERE repo_name = $1 AND user_id = $2",
		nameRepo, userId)

	_, err := result.RowsAffected()
	if err != nil {
		log.Error(err.Error())
		return custom_error.DelBookmarkFail
	}

	return nil
}
