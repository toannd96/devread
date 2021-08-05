package crawler

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/gocolly/colly/v2"
	"go.uber.org/zap"

	"devread/custom_error"
	"devread/helper"
	"devread/log"
	"devread/model"
	"devread/repository"
)

func YellowcodePost(postRepo repository.PostRepo) {
	c := colly.NewCollector()
	log := log.WriteLog()

	posts := []model.Post{}
	var yellowcodePost model.Post
	c.OnHTML("header[class=entry-header]", func(e *colly.HTMLElement) {
		yellowcodePost.Name = e.ChildText("h2.entry-title > a")
		yellowcodePost.Link = e.ChildAttr("h2.entry-title > a", "href")
		yellowcodePost.Tag = strings.ToLower(strings.Replace(
			strings.Replace(
				strings.Replace(
					e.ChildText("span.meta-category > a"), "\n", "", -1), "/", "", -1), "-", "", -1))
		yellowcodePost.PostID = helper.Hash(yellowcodePost.Name, yellowcodePost.Link)
		posts = append(posts, yellowcodePost)
	})

	c.OnScraped(func(r *colly.Response) {
		queue := helper.NewJobQueue(runtime.NumCPU())
		queue.Start()
		defer queue.Stop()

		for _, post := range posts {
			queue.Submit(&YellowcodeProcess{
				post:     post,
				postRepo: postRepo,
			})
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Error("Lỗi: ", zap.String("Truy cập ", r.Request.URL.String()), zap.Error(err))
	})

	listURL := []string{}
	for numb := 1; numb < 7; numb++ {
		trend := fmt.Sprintf("https://yellowcodebooks.com/category/lap-trinh-android/page/%d", numb)
		listURL = append(listURL, trend)
	}
	for numb := 1; numb < 6; numb++ {
		newest := fmt.Sprintf("https://yellowcodebooks.com/category/lap-trinh-java/page/%d", numb)
		listURL = append(listURL, newest)
	}

	for _, url := range listURL {
		log.Sugar().Info("Truy cập: ", url)
		c.Visit(url)
	}
}

type YellowcodeProcess struct {
	post     model.Post
	postRepo repository.PostRepo
}

func (process *YellowcodeProcess) Process() {
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
