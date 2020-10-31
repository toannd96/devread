package helper

import (
	"context"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/labstack/gommon/log"
	"runtime"
	"strings"
	"tech_posts_trending/custom_error"
	"tech_posts_trending/model"
	"tech_posts_trending/repository"
	"time"
)

func VibloPost(postRepo repository.PostRepo) {
	c := colly.NewCollector()
	c.SetRequestTimeout(30 * time.Second)

	posts := []model.Post{}
	var vibloPost model.Post
	c.OnHTML("div[class=post-title--inline]", func(e *colly.HTMLElement) {

		vibloPost.Name = e.ChildText("h3.word-break > a")
		if vibloPost.Name == "" || vibloPost.Link == "https://viblo.asia" {
			return
		}
		vibloPost.Link = "https://viblo.asia" + e.ChildAttr("h3.word-break > a", "href")
		vibloPost.Tags = strings.Replace(strings.Replace(e.ChildText("div.tags > a:last-child"), "\n", "", -1), "Trending", "", -1)
		// convert string tags to slice
		//tags := strings.Replace(strings.Replace(e.ChildText("div.tags > a"), "\n", "", -1), "Trending", "", -1)
		//vibloPost.Tags = strings.Fields(tags)

		posts = append(posts, vibloPost)
	})

	c.OnHTML(".series-header .series-title-box", func(e *colly.HTMLElement) {
		vibloPost.Name = e.ChildText("h1.series-title-header  > a")
		if vibloPost.Name == "" || vibloPost.Link == "https://viblo.asia" {
			return
		}
		vibloPost.Link = "https://viblo.asia" + e.ChildAttr("h1.series-title-header  > a", "href")
		vibloPost.Tags = e.ChildText("div.tags > a:last-child")
		posts = append(posts, vibloPost)
	})

	c.OnScraped(func(r *colly.Response) {
		queue := NewJobQueue(runtime.NumCPU())
		queue.Start()
		defer queue.Stop()

		for _, post := range posts {
			queue.Submit(&VibloProcess{
				post:       post,
				postRepo:   postRepo,
			})
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Error("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
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
	for _,url := range listURL {
		c.Visit(url)
		fmt.Println(url)
	}
}

type VibloProcess struct {
	post       model.Post
	postRepo  repository.PostRepo
}

func (process *VibloProcess) Process() {
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
