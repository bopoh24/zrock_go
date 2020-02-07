package pgstore

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
)

// TestDB helper testing function for store
func TestDB(t *testing.T, datebaseURL string) (*sql.DB, func(...string)) {
	t.Helper()

	db, err := sql.Open("postgres", datebaseURL)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	return db, func(tables ...string) {
		if len(tables) > 0 {
			db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", ")))
		}
		db.Close()
	}
}
