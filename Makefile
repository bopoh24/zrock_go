.PHONY: run
run: migrate swag build 
migrate:
	migrate -path migrations -database postgres://postgres@localhost/zrock_api_dev?sslmode=disable up
swag:
	swag init -g ./cmd/zrock_api/main.go -o ./api
build:
	go build -v ./cmd/zrock_api

.PHONY: test
test: migrations_test execute_tests
migrations_test:
	migrate -path migrations -database postgres://postgres@localhost/zrock_api_test?sslmode=disable up

execute_tests:
	go test -v -race -timeout 30s ./...



.DEFAULT_GOAL := run

