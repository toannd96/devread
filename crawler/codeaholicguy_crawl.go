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
	"time"

	"github.com/gocolly/colly/v2"
	"go.uber.org/zap"
)

func CodeaholicguyPost(postRepo repository.PostRepo) {
	log := log.WriteLog()

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
		codeaholicguyPost.PostID = helper.Hash(codeaholicguyPost.Name, codeaholicguyPost.Link)
		c.Visit(e.Request.AbsoluteURL(codeaholicguyPost.Link))
		if codeaholicguyPost.Name == "" || codeaholicguyPost.Link == "" {
			return
		}
		posts = append(posts, codeaholicguyPost)
	})

	c.OnScraped(func(r *colly.Response) {
		queue := helper.NewJobQueue(runtime.NumCPU())
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
}

func (process *CodeaholicguyProcess) Process() {
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
