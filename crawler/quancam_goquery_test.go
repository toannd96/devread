package crawler

import (
	"context"
	"devread/custom_error"
	"devread/helper"
	"devread/model"
	"devread/repository"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
	"log"
	"runtime"
	"strconv"
	"strings"
)

const urlBaseTest = "https://quan-cam.com"

func checkErrorTest(err error) {
	if err != nil {
		log.Println(err)
	}
}

func GetListPage() []string {
	pageList := make([]string, 0)
	page := []int{1}
	for len(page) > 0 {
		pathURL := fmt.Sprintf("https://quan-cam.com/posts?page=%d", page[0])
		fmt.Println("GetListPage")
		response, err := helper.HttpClient.GetRequestWithRetries(pathURL)
		if err != nil {
			log.Println(err)
		}
		defer response.Body.Close()
		doc, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			log.Println(err)
		}

		link, _ := doc.Find("a.next").Attr("href")
		if link != "" {
			split := strings.Split(link, "=")[1]
			nextLink, _ := strconv.Atoi(split)
			page[0] = nextLink
			url := fmt.Sprintf("https://quan-cam.com/posts?page=%d", nextLink)
			pageList = append(pageList, url)
		} else {
			page = page[:0]
		}
	}
	fmt.Println("list page->", pageList)
	return pageList
}

func getOnePageTest(pathURL string) ([]model.Post, error) {
	response, err := helper.HttpClient.GetRequestWithRetries(pathURL)
	checkError(err)
	defer response.Body.Close()
	doc, err := goquery.NewDocumentFromReader(response.Body)
	checkError(err)

	posts := make([]model.Post, 0)
	doc.Find("div[class=post]").Each(func(i int, s *goquery.Selection) {
		var quancamPost model.Post
		quancamPost.Name = s.Find("h3.post__title > a").Text()
		link, _ := s.Find("h3.post__title > a").Attr("href")
		quancamPost.Link = urlBase + link
		quancamPost.Tag = strings.ToLower(strings.Replace(
			strings.Replace(
				s.Find("span.tagging > a").Text(), "\n", "", -1), "#", " ", -1))
		posts = append(posts, quancamPost)

		//fmt.Println("Name", quancamPost.Name)
		//fmt.Println("Link", quancamPost.Link)
		//fmt.Println("Tag", quancamPost.Tag)
		//fmt.Println("\n ")
	})
	return posts, nil
}

func QuancamPostV3(postRepo repository.PostRepo) {
	sem := semaphore.NewWeighted(int64(runtime.NumCPU()))
	group, ctx := errgroup.WithContext(context.Background())
	listPage := GetListPage()

	for _, page := range listPage {
		err := sem.Acquire(ctx, 1)
		if err != nil {
			fmt.Printf("Acquire err = %+v\n", err)
			continue
		}
		group.Go(func() error {
			defer sem.Release(1)

			//do work
			posts, err := getOnePage(page)
			checkError(err)
			queue := helper.NewJobQueue(runtime.NumCPU())
			queue.Start()
			defer queue.Stop()
			for _, post := range posts {
				queue.Submit(&QuancamProcessV3{
					post:     post,
					postRepo: postRepo,
				})
			}
			return nil
		})
	}
	if err := group.Wait(); err != nil {
		fmt.Printf("g.Wait() err = %+v\n", err)
	}
}

type QuancamProcessV3 struct {
	post     model.Post
	postRepo repository.PostRepo
}

func (process *QuancamProcessV3) Process() {
	// select post by name
	cacheRepo, err := process.postRepo.SelectPostByName(context.Background(), process.post.Name)
	if err == custom_error.PostNotFound {
		// insert post to database
		fmt.Println("Add: ", process.post.Name)
		_, err = process.postRepo.SavePost(context.Background(), process.post)
		checkError(err)
		return
	}

	// update post
	if process.post.Name != cacheRepo.Name ||
		process.post.Link != cacheRepo.Link {
		fmt.Println("Updated: ", process.post.Name)
		_, err = process.postRepo.UpdatePost(context.Background(), process.post)
		checkError(err)
	}
}
