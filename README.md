## DevRead - ứng dụng tổng hợp kiến thức cho developer
- Tổng hợp bài viết hay nhất trên các blog IT như viblo, toidicodedao, yellowcodebooks, thefullsnack, quan-cam, codeaholicguy,...
- Nội dung thu thập như trong [tệp](https://github.com/dactoankmapydev/devread/blob/master/huong_dan/posts.csv) không vi phạm bản quyền tác giả
- Cấu trúc:
```
.
├── 1_init.sql
├── crawler
│   ├── codeaholicguy_crawl.go
│   ├── quancam_v1.go
│   ├── quancam_v2.go
│   ├── thefullsnack_crawl.go
│   ├── toidicodedao_crawl.go
│   ├── viblo_crawl.go
│   └── yellowcode_crawl.go
├── custom_error
│   ├── err_post.go
│   └── err_user.go
├── db
│   ├── db.go
│   └── redis.go
├── docker-compose.yml
├── Dockerfile
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── handler
│   ├── post_handler.go
│   └── user_handler.go
├── helper
│   ├── custom_validator.go
│   ├── hash_value.go
│   ├── hash_values.go
│   ├── http_client.go
│   └── job_queue.go
├── log
│   └── log.go
├── log_files
│   └── error
│       └── devread_20201224_error.log
├── main.go
├── middleware
│   ├── gzip.go
│   ├── jwt.go
│   └── req_headers.go
├── model
│   ├── post.go
│   ├── req
│   │   ├── pwd_submit.go
│   │   ├── req_bookmark.go
│   │   ├── req_email.go
│   │   ├── req_signin.go
│   │   ├── req_signup.go
│   │   ├── req_tag.go
│   │   └── req_user_update.go
│   ├── resp.go
│   ├── token.go
│   └── user.go
├── README.md
├── repository
│   ├── authen_repo.go
│   ├── bookmark_repo.go
│   ├── post_repo.go
│   ├── repo_impl
│   │   ├── authen_repo_impl.go
│   │   ├── bookmark_repo_impl.go
│   │   ├── post_repo_impl.go
│   │   └── user_repo_impl.go
│   └── user_repo.go
├── router
│   └── api.go
└── security
    ├── jwt.go
    ├── password.go
    └── token.go
```

- Danh sách API hiện có:

![](https://github.com/dactoankmapydev/devread/blob/master/huong_dan/api.png)

- Danh sách bài viết:

![](https://github.com/dactoankmapydev/devread/blob/master/huong_dan/posts.jpg)

- Hồ sơ cá nhân:

![](https://github.com/dactoankmapydev/devread/blob/master/huong_dan/profile.jpg)

- Bộ sưu tập:

![](https://github.com/dactoankmapydev/devread/blob/master/huong_dan/collection.jpg)

- Đăng ký:

![](https://github.com/dactoankmapydev/devread/blob/master/huong_dan/signup.jpg)

- Đăng nhập:

![](https://github.com/dactoankmapydev/devread/blob/master/huong_dan/signin.jpg)

## Run
```
docker-compose up -d
docker-compose up
```

- Run swagger test API tại ```localhost:3000/swagger/index.html```
