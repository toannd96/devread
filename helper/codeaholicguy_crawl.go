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

func CodeaholicguyPost(postRepo repository.PostRepo) {
	c := colly.NewCollector()
	c.SetRequestTimeout(30 * time.Second)

	posts := []model.Post{}
	var codeaholicguyPost model.Post

	c.OnHTML("span[class=cat-links]", func(e *colly.HTMLElement) {
		if codeaholicguyPost.Name == "" || codeaholicguyPost.Link == "" {
			return
		}
		codeaholicguyPost.Tags = strings.ToLower(strings.Replace(e.ChildText("span.cat-links > a:last-child"), "Chuyện coding", "", -1))
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
		queue := NewJobQueue(runtime.NumCPU())
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
		log.Error("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	for i := 1; i < 7; i++ {
		fullURL := fmt.Sprintf("https://codeaholicguy.com/category/chuyen-coding/page/%d", i)
		c.Visit(fullURL)
		fmt.Println(fullURL)
	}
}

type CodeaholicguyProcess struct {
	post     model.Post
	postRepo repository.PostRepo
}

func (process *CodeaholicguyProcess) Process() {
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
