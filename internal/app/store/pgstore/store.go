package pgstore

import (
	"database/sql"

	"github.com/bopoh24/zrock_go/internal/app/store"
	_ "github.com/lib/pq" // postgres driver
)

// Store struct for store
type Store struct {
	db             *sql.DB
	userRepository *UserRepository
}

// New creates Store
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// User is a UserRepository instance
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}
	s.userRepository = &UserRepository{
		store: s,
	}
	return s.userRepository
}
