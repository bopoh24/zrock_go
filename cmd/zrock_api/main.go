package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/bopoh24/zrock_go/internal/app/apiserver"
	"github.com/bopoh24/zrock_go/internal/app/settings"
	"github.com/bopoh24/zrock_go/internal/app/store"
	"github.com/bopoh24/zrock_go/internal/app/store/pgstore"
)

func initDbStore() (store.IfaceStore, *sql.DB, error) {
	db, err := sql.Open("postgres", settings.App.DatabaseURL)
	if err != nil {
		return nil, nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, nil, err
	}
	return pgstore.New(db), db, nil
}

func main() {
	store, db, err := initDbStore()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	srv := apiserver.NewServer(store)
	http.ListenAndServe(settings.App.BindAdd, srv)
}
