package crawler

import (
	"devread/custom_error"
	"devread/helper"
	"devread/log"
	"devread/model"
	"devread/repository"

	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

const urlBase = "https://quan-cam.com"

func getOnePage(pathURL string) ([]model.Post, error) {
	log := log.WriteLog()
	response, err := helper.GetRequestWithRetries(pathURL)
	if err != nil {
		log.Error("Lỗi: ", zap.Error(err))
	}

	defer response.Body.Close()
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Error("Lỗi: ", zap.Error(err))
	}

	posts := make([]model.Post, 0)
	doc.Find("div[class=post]").Each(func(i int, s *goquery.Selection) {
		var quancamPost model.Post
		quancamPost.Name = s.Find("h3.post__title > a").Text()
		link, _ := s.Find("h3.post__title > a").Attr("href")
		quancamPost.Link = urlBase + link
		quancamPost.Tag = strings.ToLower(strings.Replace(
			strings.Replace(
				s.Find("span.tagging > a").Text(), "\n", "", -1), "#", " ", -1))
		quancamPost.PostID = helper.Hash(quancamPost.Name, quancamPost.Link)
		posts = append(posts, quancamPost)
	})
	return posts, nil
}

func QuancamPostV1(postRepo repository.PostRepo) {
	log := log.WriteLog()

	sem := semaphore.NewWeighted(int64(runtime.NumCPU()))
	group, ctx := errgroup.WithContext(context.Background())

	for page := 1; page <= 4; page++ {
		pathURL := fmt.Sprintf("%s/posts?page=%d", urlBase, page)
		err := sem.Acquire(ctx, 1)
		if err != nil {
			continue
		}
		group.Go(func() error {
			defer sem.Release(1)

			//do work
			posts, err := getOnePage(pathURL)
			if err != nil {
				log.Error("Lỗi: ", zap.Error(err))
			}

			queue := helper.NewJobQueue(runtime.NumCPU())
			queue.Start()
			defer queue.Stop()
			for _, post := range posts {
				queue.Submit(&QuancamProcessV1{
					post:     post,
					postRepo: postRepo,
				})
			}

			return nil
		})
	}
	if err := group.Wait(); err != nil {
		log.Error("Có 1 goroutine lỗi ", zap.Error(err))
	}
}

type QuancamProcessV1 struct {
	post     model.Post
	postRepo repository.PostRepo
}

func (process *QuancamProcessV1) Process() {
	log := log.WriteLog()
	// select post by id
	cacheRepo, err := process.postRepo.SelectById(context.Background(), process.post.PostID)
	if err == custom_error.PostNotFound {
		// insert post to database
		log.Sugar().Info("Thêm bài viết: ", process.post.Name)
		_, err = process.postRepo.Save(context.Background(), process.post)
		if err != nil {
			log.Error("Thêm bài viết thất bại ", zap.Error(err))
		}
		return
	}

	// update post
	if process.post.PostID != cacheRepo.PostID {
		log.Sugar().Info("Thêm bài viết: ", process.post.Name)
		_, err = process.postRepo.Update(context.Background(), process.post)
		if err != nil {
			log.Error("Thêm bài viết thất bại ", zap.Error(err))
		}
	}
}
