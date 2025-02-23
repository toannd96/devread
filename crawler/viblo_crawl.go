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

func VibloPost(postRepo repository.PostRepo) {
	c := colly.NewCollector()
	log, _ := handle_log.WriteLog()

	c.SetRequestTimeout(30 * time.Second)

	posts := []model.Post{}
	var vibloPost model.Post
	c.OnHTML("div[class=post-title--inline]", func(e *colly.HTMLElement) {

		vibloPost.Name = e.ChildText("h3.word-break > a")
		if vibloPost.Name == "" || vibloPost.Link == "https://viblo.asia" {
			return
		}
		vibloPost.Link = "https://viblo.asia" + e.ChildAttr("h3.word-break > a", "href")
		vibloPost.Tag = strings.ToLower(
			strings.Replace(
				strings.Replace(e.ChildText("div.tags > a:last-child"), "\n", "", -1), "Trending", "", -1))
		posts = append(posts, vibloPost)
	})

	c.OnHTML(".series-header .series-title-box", func(e *colly.HTMLElement) {
		vibloPost.Name = e.ChildText("h1.series-title-header  > a")
		if vibloPost.Name == "" || vibloPost.Link == "https://viblo.asia" {
			return
		}
		vibloPost.Link = "https://viblo.asia" + e.ChildAttr("h1.series-title-header  > a", "href")
		vibloPost.Tag = strings.ToLower(e.ChildText("div.tags > a:last-child"))
		posts = append(posts, vibloPost)
	})

	c.OnScraped(func(r *colly.Response) {
		queue := helper.NewJobQueue(2)
		queue.Start()
		defer queue.Stop()

		for _, post := range posts {
			queue.Submit(&VibloProcess{
				post:     post,
				postRepo: postRepo,
			})
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Error("Lỗi: ", zap.String("Truy cập ", r.Request.URL.String()), zap.Error(err))
	})

	listURL := []string{}
	for numb := 1; numb < 5; numb++ {
		trend := fmt.Sprintf("https://viblo.asia/trending?page=%d", numb)
		listURL = append(listURL, trend)
	}
	for numb := 1; numb < 4; numb++ {
		newest := fmt.Sprintf("https://viblo.asia/newest?page=%d", numb)
		listURL = append(listURL, newest)
	}
	for numb := 1; numb < 34; numb++ {
		series := fmt.Sprintf("https://viblo.asia/series?page=%d", numb)
		listURL = append(listURL, series)
	}
	for _, url := range listURL {
		log.Sugar().Info("Truy cập: ", url)
		c.Visit(url)
	}
}

type VibloProcess struct {
	post     model.Post
	postRepo repository.PostRepo
	logger   *zap.Logger
}

func (process *VibloProcess) Process() {
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
