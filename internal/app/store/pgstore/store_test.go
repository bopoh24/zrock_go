package pgstore

import (
	"os"
	"testing"
)

var databaseURL string

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "host=localhost dbname=zrock_api_test sslmode=disable user=postgres"
	}
	os.Exit(m.Run())
}
