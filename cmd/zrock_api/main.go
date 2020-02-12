package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/bopoh24/zrock_go/internal/app/apiserver"
	"github.com/bopoh24/zrock_go/internal/app/settings"
	"github.com/bopoh24/zrock_go/internal/app/store/pgstore"
)

func initDbStore() (*pgstore.Store, *sql.DB, error) {
	db, err := sql.Open("postgres", settings.App.DatabaseURL)
	if err != nil {
		return nil, nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, nil, err
	}
	return pgstore.New(db), db, nil
}

// @title Zrock API
// @version 1.0
// @description Zrock REST API Server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	store, db, err := initDbStore()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	srv := apiserver.NewServer(store)

	http.ListenAndServe(settings.App.BindAdd, srv)
}
