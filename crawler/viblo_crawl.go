package crawler

import (
	"context"
	"devread/custom_error"
	"devread/helper"
	"devread/model"
	"devread/repository"
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"runtime"
	"strings"
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
		vibloPost.PostID = helper.Hash(vibloPost.Name, vibloPost.Link)
		posts = append(posts, vibloPost)
	})

	c.OnScraped(func(r *colly.Response) {
		queue := helper.NewJobQueue(runtime.NumCPU())
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
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
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
		c.Visit(url)
		fmt.Println(url)
	}
}

type VibloProcess struct {
	post     model.Post
	postRepo repository.PostRepo
}

func (process *VibloProcess) Process() {
	// select post by id
	cacheRepo, err := process.postRepo.SelectById(context.Background(), process.post.PostID)
	if err == custom_error.PostNotFound {
		// insert post to database
		fmt.Println("Add: ", process.post.Name)
		_, err = process.postRepo.Save(context.Background(), process.post)
		if err != nil {
			log.Println(err)
		}
		return
	}

	// update post
	if process.post.PostID != cacheRepo.PostID {
		fmt.Println("Updated: ", process.post.Name)
		_, err = process.postRepo.Update(context.Background(), process.post)
		if err != nil {
			log.Println(err)
		}
	}
}
