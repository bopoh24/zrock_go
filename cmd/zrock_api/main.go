package main

import (
	"net/http"

	"github.com/bopoh24/zrock_go/internal/app/zrockapi"
)

func main() {
	config := zrockapi.NewConfig()

	srv := zrockapi.NewServer(config)

	http.ListenAndServe(config.BindAdd, srv)
}
