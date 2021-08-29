package crawler

import (
	"devread/custom_error"
	"devread/handle_log"
	"devread/helper"
	"devread/model"
	"devread/repository"

	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"go.uber.org/zap"
)

func CodeaholicguyPost(postRepo repository.PostRepo) {
	log, _ := handle_log.WriteLog()

	c := colly.NewCollector()
	c.SetRequestTimeout(30 * time.Second)

	posts := []model.Post{}
	var codeaholicguyPost model.Post

	c.OnHTML("span[class=cat-links]", func(e *colly.HTMLElement) {
		if codeaholicguyPost.Name == "" || codeaholicguyPost.Link == "" {
			return
		}
		codeaholicguyPost.Tag = strings.ToLower(strings.Replace(e.ChildText("span.cat-links > a:last-child"), "Chuyện coding", "", -1))
		posts = append(posts, codeaholicguyPost)
	})

	c.OnHTML("header[class=entry-header]", func(e *colly.HTMLElement) {
		codeaholicguyPost.Name = e.ChildText("h1.entry-title > a")
		codeaholicguyPost.Link = e.ChildAttr("h1.entry-title > a", "href")
		c.Visit(e.Request.AbsoluteURL(codeaholicguyPost.Link))
		if codeaholicguyPost.Name == "" || codeaholicguyPost.Link == "" {
			return
		}
		posts = append(posts, codeaholicguyPost)
	})

	c.OnScraped(func(r *colly.Response) {
		queue := helper.NewJobQueue(2)
		queue.Start()
		defer queue.Stop()

		for _, post := range posts {
			queue.Submit(&CodeaholicguyProcess{
				post:     post,
				postRepo: postRepo,
			})
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Error("Lỗi: ", zap.String("Truy cập ", r.Request.URL.String()), zap.Error(err))
	})
	for i := 1; i < 7; i++ {
		fullURL := fmt.Sprintf("https://codeaholicguy.com/category/chuyen-coding/page/%d", i)
		log.Sugar().Info("Truy cập: ", fullURL)
		c.Visit(fullURL)
	}
}

type CodeaholicguyProcess struct {
	post     model.Post
	postRepo repository.PostRepo
	logger   *zap.Logger
}

func (process *CodeaholicguyProcess) Process() {
	if process.logger == nil {
		l, _ := handle_log.WriteLog()
		process.logger = l
	}

	// select post by link
	cacheRepo, err := process.postRepo.SelectByLink(context.Background(), process.post.Link)
	if err == custom_error.PostNotFound {
		// insert post to database
		process.logger.Sugar().Info("Thêm bài viết: ", process.post.Name)
		_, err = process.postRepo.Save(context.Background(), process.post)
		if err != nil {
			process.logger.Error("Thêm bài viết thất bại ", zap.String("bài viết: ", process.post.Name), zap.Error(err))
		}
		return
	}

	// update post
	if process.post.Name != cacheRepo.Name {
		process.logger.Sugar().Info("Cập nhật bài viết: ", process.post.Name)
		_, err = process.postRepo.Update(context.Background(), process.post)
		if err != nil {
			process.logger.Error("Cập nhật bài viết thất bại ", zap.String("bài viết: ", process.post.Name), zap.Error(err))
		}
	}
}
