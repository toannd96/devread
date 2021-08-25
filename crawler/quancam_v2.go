package crawler

import (
	"devread/custom_error"
	"devread/handle_log"
	"devread/helper"
	"devread/model"
	"devread/repository"

	"context"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

const urlBaseTest = "https://quan-cam.com"

func GetListPage() []string {
	log, _ := handle_log.WriteLog()

	pageList := make([]string, 0)
	page := []int{1}
	for len(page) > 0 {
		pathURL := fmt.Sprintf("https://quan-cam.com/posts?page=%d", page[0])
		response, err := helper.GetRequestWithRetries(pathURL)
		if err != nil {
			log.Error("Lỗi: ", zap.Error(err))
		}

		defer response.Body.Close()
		doc, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			log.Error("Lỗi: ", zap.Error(err))
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
	log.Sugar().Info("Danh sách trang ", pageList)
	return pageList
}

func getOnePageTest(pathURL string) ([]model.Post, error) {
	log, _ := handle_log.WriteLog()

	response, err := helper.GetRequestWithRetries(pathURL)
	if err != nil {
		log.Error("Lỗi: ", zap.Error(err))
	}

	defer response.Body.Close()
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Error("Lỗi: ", zap.Error(err))
	}

	posts := make([]model.Post, 0)
	doc.Find("div[class=post]").Each(func(i int, s *goquery.Selection) {
		var quancamPost model.Post
		quancamPost.Name = s.Find("h3.post__title > a").Text()
		link, _ := s.Find("h3.post__title > a").Attr("href")
		quancamPost.Link = urlBase + link
		quancamPost.Tag = strings.ToLower(strings.Replace(
			strings.Replace(
				s.Find("span.tagging > a").Text(), "\n", "", -1), "#", " ", -1))
		quancamPost.PostID = helper.Hash(quancamPost.Name, quancamPost.Link)
		posts = append(posts, quancamPost)

		log.Sugar().Info("Tên bài viết ", quancamPost.Name)
		log.Sugar().Info("Link", quancamPost.Link)
		log.Sugar().Info("Tag", quancamPost.Tag)
	})
	return posts, nil
}

func QuancamPostV2(postRepo repository.PostRepo) {
	log, _ := handle_log.WriteLog()

	sem := semaphore.NewWeighted(int64(runtime.NumCPU()))
	group, ctx := errgroup.WithContext(context.Background())
	listPage := GetListPage()

	for _, page := range listPage {
		err := sem.Acquire(ctx, 1)
		if err != nil {
			continue
		}
		group.Go(func() error {
			defer sem.Release(1)

			//do work
			posts, err := getOnePageTest(page)
			if err != nil {
				log.Error("Lỗi: ", zap.Error(err))
			}

			queue := helper.NewJobQueue(runtime.NumCPU())
			queue.Start()
			defer queue.Stop()
			for _, post := range posts {
				queue.Submit(&QuancamProcessV2{
					post:     post,
					postRepo: postRepo,
				})
			}
			return nil
		})
	}
	if err := group.Wait(); err != nil {
		log.Error("Có 1 goroutine lỗi ", zap.Error(err))
	}
}

type QuancamProcessV2 struct {
	post     model.Post
	postRepo repository.PostRepo
	logger   *zap.Logger
}

func (process *QuancamProcessV2) Process() {
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
