package main

import (
	"net/http"

	"github.com/bopoh24/zrock_go/internal/app/zrock_api"
)

func main() {
	srv := zrock_api.NewServer()
	http.ListenAndServe(":8080", srv)
}
