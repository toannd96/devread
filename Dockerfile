# Start from golang base image
FROM golang:latest as builder

# Install git.
# Git is required for fetching the dependencies.
RUN apt-get update && apt install git

ENV GO111MODULE=on
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/app" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

# Set the current working directory inside the container 
WORKDIR $GOPATH/app

# Copy go mod and sum files 
COPY go.mod go.sum ./

# Copy the source from the current directory to the working Directory inside the container 
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o /main .

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed 
RUN go mod download

# Expose port to the outside world
EXPOSE 3000

#Command to run the executable
CMD ./main
