package teststore

import (
	"github.com/bopoh24/zrock_go/internal/app/model"
	"github.com/bopoh24/zrock_go/internal/app/store"
)

type UserRepository struct {
	store *Store
	users map[int]*model.User
}

func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}
	if err := u.BeforeCreate(); err != nil {
		return err
	}
	u.ID = len(r.users) + 1
	r.users[u.ID] = u
	return nil
}

func (r *UserRepository) FindByPk(userID int) (*model.User, error) {
	u, ok := r.users[userID]
	if !ok {
		return nil, store.ErrRecordNotFound
	}
	return u, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, store.ErrRecordNotFound
}

func (r *UserRepository) FindByNick(nick string) (*model.User, error) {
	for _, u := range r.users {
		if u.Nick == nick {
			return u, nil
		}
	}
	return nil, store.ErrRecordNotFound
}
