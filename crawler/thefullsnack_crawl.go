package crawler

import (
	"devread/custom_error"
	"devread/helper"
	"devread/log"
	"devread/model"
	"devread/repository"

	"github.com/gocolly/colly/v2"
	"go.uber.org/zap"

	"context"
	"regexp"
	"runtime"
	"strings"
)

func ThefullsnackPost(postRepo repository.PostRepo) {
	c := colly.NewCollector()
	log := log.WriteLog()

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
		thefullsnackPost.PostID = helper.Hash(thefullsnackPost.Name, thefullsnackPost.Link)
		posts = append(posts, thefullsnackPost)
	})

	c.OnScraped(func(r *colly.Response) {
		queue := helper.NewJobQueue(runtime.NumCPU())
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
}

func (process *ThefullsnackProcess) Process() {
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
