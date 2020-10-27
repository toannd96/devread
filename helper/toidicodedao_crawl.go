package helper

import (
	"context"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/labstack/gommon/log"
	"runtime"
	"tech_posts_trending/custom_error"
	"tech_posts_trending/model"
	"tech_posts_trending/repository"
)

func ToidicodedaoPost(postRepo repository.PostRepo) {
	c := colly.NewCollector()

	posts := make([]model.Post, 0, 11)
	c.OnHTML(".site-content .entry-title", func(e *colly.HTMLElement) {
		var toidicodedaoPost model.Post
		toidicodedaoPost.Name = e.Text
		toidicodedaoPost.Link = e.ChildAttr("h1.entry-title > a", "href")
		posts = append(posts, toidicodedaoPost)
	})

	c.OnScraped(func(r *colly.Response) {
		queue := NewJobQueue(runtime.NumCPU())
		queue.Start()
		defer queue.Stop()

		for _, post := range posts {
			queue.Submit(&ToidicodedaoProcess{
				post:       post,
				postRepo:   postRepo,
			})
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	for i := 1; i < 32; i++ {
		fullURL := fmt.Sprintf("https://toidicodedao.com/category/chuyen-coding/page/%d", i)
		c.Visit(fullURL)
		fmt.Println(fullURL)
	}
}

type ToidicodedaoProcess struct {
	post       model.Post
	postRepo  repository.PostRepo
}

func (process *ToidicodedaoProcess) Process() {
	// select post by name
	cacheRepo, err := process.postRepo.SelectPostByName(context.Background(), process.post.Name)
	if err == custom_error.PostNotFound {
		// insert post to database
		fmt.Println("Add: ", process.post.Name)
		_, err = process.postRepo.SavePost(context.Background(), process.post)
		if err != nil {
			log.Error(err)
		}
		return
	}

	// update post
	if process.post.Name != cacheRepo.Name ||
		process.post.Link != cacheRepo.Link {
		fmt.Println("Updated: ", process.post.Name)
		_, err = process.postRepo.UpdatePost(context.Background(), process.post)
		if err != nil {
			log.Error(err)
		}
	}
}
