.PHONY: build
build:
	go build -v ./cmd/zrock_api


.PHONY: test
test:
	go test -v -race -timeout 30s ./...
.DEFAULT_GOAL := build

