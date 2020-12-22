FROM golang:latest

RUN apk update && apk add git

ENV GO111MODULE=on

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/app" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH/app

COPY . .

RUN go mod init tech_posts_trending

WORKDIR /app
RUN GOOS=linux go build -o app

ENTRYPOINT ["./app"]

EXPOSE 3000