GOPATH:=$(shell go env GOPATH)

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -o goChat server.go

.PHONY: docker
docker:
	docker build . -t go-chat:v0.0.1

.PHONY: run
run:
	go mod download
	go run server.go



