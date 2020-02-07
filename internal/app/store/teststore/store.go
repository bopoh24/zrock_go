package teststore

import (
	"github.com/bopoh24/zrock_go/internal/app/model"
	"github.com/bopoh24/zrock_go/internal/app/store"
)

// Store struct for store
type Store struct {
	userRepository *UserRepository
}

// New creates Store
func New() *Store {
	return &Store{}
}

// User is a UserRepository instance
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}
	s.userRepository = &UserRepository{
		store: s,
		users: make(map[int]*model.User),
	}
	return s.userRepository
}
