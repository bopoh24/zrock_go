.PHONY: run
run: migrate swag build 
migrate:
	migrate -path migrations -database postgres://postgres@localhost/zrock_api_dev?sslmode=disable up
swag:
	swag init -g ./cmd/zrock_api/main.go -o ./api
build:
	go build -v ./cmd/zrock_api

.PHONY: test
test: migrate_test run_tests
migrate_test:
	migrate -path migrations -database postgres://postgres@localhost/zrock_api_test?sslmode=disable up

run_tests:
	go test -v -race -timeout 30s ./...


.PHONY: coverage
coverage:
	go test ./... -cover -coverprofile=./coverage.out  && go tool cover -html=./coverage.out


.DEFAULT_GOAL := run

 