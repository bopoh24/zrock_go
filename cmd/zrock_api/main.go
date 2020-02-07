package main

import (
	"net/http"

	"github.com/bopoh24/zrock_go/internal/app/apiserver"
)

func main() {
	config := apiserver.NewConfig()

	srv := apiserver.NewServer(config)

	http.ListenAndServe(config.BindAdd, srv)
}
