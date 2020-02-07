package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/bopoh24/zrock_go/internal/app/apiserver"
	"github.com/bopoh24/zrock_go/internal/app/store"
	"github.com/bopoh24/zrock_go/internal/app/store/pgstore"
)

func initDbStore(datebaseURL string) (store.IfaceStore, *sql.DB, error) {
	db, err := sql.Open("postgres", datebaseURL)
	if err != nil {
		return nil, nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, nil, err
	}

	return pgstore.New(db), db, nil
}

func main() {
	config := apiserver.NewConfig()
	store, db, err := initDbStore(config.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	srv := apiserver.NewServer(config, store)

	http.ListenAndServe(config.BindAdd, srv)
}
