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
)

func QuancamPost(postRepo repository.PostRepo) {
	c := colly.NewCollector()

	posts := []model.Post{}
	c.OnHTML("div[class=post]", func(e *colly.HTMLElement) {
		var quancamPost model.Post
		quancamPost.Name = e.ChildText("h3.post__title > a")
		quancamPost.Link = "https://quan-cam.com" + e.ChildAttr("h3.post__title > a", "href")
		quancamPost.Tag = strings.ToLower(strings.Replace(
			strings.Replace(
				e.ChildText("span.tagging > a"), "\n", "", -1), "#", " ", -1))
		posts = append(posts, quancamPost)
	})

	c.OnScraped(func(r *colly.Response) {
		queue := NewJobQueue(runtime.NumCPU())
		queue.Start()
		defer queue.Stop()

		for _, post := range posts {
			queue.Submit(&QuancamProcess{
				post:     post,
				postRepo: postRepo,
			})
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Error("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	for i := 1; i < 5; i++ {
		fullURL := fmt.Sprintf("https://quan-cam.com/posts?page=%d", i)
		c.Visit(fullURL)
		fmt.Println(fullURL)
	}
}

type QuancamProcess struct {
	post     model.Post
	postRepo repository.PostRepo
}

func (process *QuancamProcess) Process() {
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
