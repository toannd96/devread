package crawler

import (
	"devread/custom_error"
	"devread/handle_log"
	"devread/helper"
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

func ToidicodedaoPost(postRepo repository.PostRepo) {
	c := colly.NewCollector()
	log, _ := handle_log.WriteLog()
	c.SetRequestTimeout(30 * time.Second)

	posts := []model.Post{}
	var toidicodedaoPost model.Post

	c.OnHTML("footer[class=entry-meta]", func(e *colly.HTMLElement) {
		if toidicodedaoPost.Name == "" || toidicodedaoPost.Link == "" {
			return
		}
		toidicodedaoPost.Tag = strings.ToLower(e.ChildText("span.tag-links > a:last-child"))
		posts = append(posts, toidicodedaoPost)
	})

	c.OnHTML(".site-content .entry-title", func(e *colly.HTMLElement) {
		toidicodedaoPost.Name = e.Text
		toidicodedaoPost.Link = e.ChildAttr("h1.entry-title > a", "href")
		c.Visit(e.Request.AbsoluteURL(toidicodedaoPost.Link))
		if toidicodedaoPost.Name == "" || toidicodedaoPost.Link == "" {
			return
		}
		posts = append(posts, toidicodedaoPost)
	})

	c.OnScraped(func(r *colly.Response) {
		queue := helper.NewJobQueue(runtime.NumCPU())
		queue.Start()
		defer queue.Stop()

		for _, post := range posts {
			queue.Submit(&ToidicodedaoProcess{
				post:     post,
				postRepo: postRepo,
			})
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Error("Lỗi: ", zap.String("Truy cập ", r.Request.URL.String()), zap.Error(err))
	})

	for i := 1; i < 32; i++ {
		fullURL := fmt.Sprintf("https://toidicodedao.com/category/chuyen-coding/page/%d", i)
		log.Sugar().Info("Truy cập: ", fullURL)
		c.Visit(fullURL)
	}
}

type ToidicodedaoProcess struct {
	post     model.Post
	postRepo repository.PostRepo
	logger   *zap.Logger
}

func (process *ToidicodedaoProcess) Process() {
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
