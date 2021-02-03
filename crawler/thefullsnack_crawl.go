package crawler

import (
	"context"
	"devread/custom_error"
	"devread/helper"
	"devread/model"
	"devread/repository"
	"github.com/gocolly/colly/v2"
	"github.com/labstack/gommon/log"
	"regexp"
	"runtime"
	"strings"
)

func ThefullsnackPost(postRepo repository.PostRepo) {
	c := colly.NewCollector()

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
		log.Error("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit("https://thefullsnack.com/")
}

type ThefullsnackProcess struct {
	post     model.Post
	postRepo repository.PostRepo
}

func (process *ThefullsnackProcess) Process() {
	// select post by name
	cacheRepo, err := process.postRepo.SelectPostByName(context.Background(), process.post.Name)
	if err == custom_error.PostNotFound {
		// insert post to database
		_, err = process.postRepo.SavePost(context.Background(), process.post)
		if err != nil {
			log.Error(err)
		}
		return
	}

	// update post
	if process.post.Name != cacheRepo.Name {
		log.Info("Updated: ", process.post.Name)
		_, err = process.postRepo.UpdatePost(context.Background(), process.post)
		if err != nil {
			log.Error(err)
		}
	}
}
