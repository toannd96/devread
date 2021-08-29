package crawler

import (
	"devread/custom_error"
	"devread/handle_log"
	"devread/helper"
	"devread/model"
	"devread/repository"

	"github.com/gocolly/colly/v2"
	"go.uber.org/zap"

	"context"
	"regexp"
	"strings"
)

func ThefullsnackPost(postRepo repository.PostRepo) {
	c := colly.NewCollector()
	log, _ := handle_log.WriteLog()

	posts := []model.Post{}
	c.OnHTML("div[class=home-list-item]", func(e *colly.HTMLElement) {
		var thefullsnackPost model.Post
		thefullsnackPost.Name = e.ChildText("div.home-list-item > a")
		thefullsnackPost.Link = "https://thefullsnack.com" + e.ChildAttr("div.home-list-item > a", "href")
		tags := strings.ToLower(e.Text)
		regexSplitName := regexp.MustCompile("[0-9]{2}[-]{1}[0-9]{2}[-]{1}[0-9]{4}([a-z]{1,60}[-][a-z]{1,60}|[a-z]{1,60}|)|[,]\\s([a-z]{1,60}[-][a-z]{1,60}|[a-z]{1,60}|)")
		regexSplitTime := regexp.MustCompile("[0-9]{2}[-]{1}[0-9]{2}[-]{1}[0-9]{4}")
		splitName := strings.Join(regexSplitName.FindAllString(tags, -1), " ")
		splitTime := strings.Join(regexSplitTime.FindAllString(splitName, -1), " ")
		thefullsnackPost.Tag = strings.Replace(splitName, splitTime, "", -1)
		posts = append(posts, thefullsnackPost)
	})

	c.OnScraped(func(r *colly.Response) {
		queue := helper.NewJobQueue(2)
		queue.Start()
		defer queue.Stop()

		for _, post := range posts {
			queue.Submit(&ThefullsnackProcess{
				post:     post,
				postRepo: postRepo,
			})
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Error("Lỗi: ", zap.String("Truy cập ", r.Request.URL.String()), zap.Error(err))
	})

	c.Visit("https://thefullsnack.com/")
}

type ThefullsnackProcess struct {
	post     model.Post
	postRepo repository.PostRepo
	logger   *zap.Logger
}

func (process *ThefullsnackProcess) Process() {
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
