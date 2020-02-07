package teststore

import (
	"log"

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
	if err := r.checkEmailAndNickFree(u); err != nil {
		return err
	}
	if err := u.BeforeCreate(); err != nil {
		return err
	}
	u.ID = len(r.users) + 1
	r.users[u.ID] = u
	return nil
}

func (r *UserRepository) checkEmailAndNickFree(u *model.User) error {
	log.Print(r.users)
	for _, user := range r.users {

		if u.Email == user.Email {
			return store.ErrEmailExists
		}
		if u.Nickname == user.Nickname {
			return store.ErrNicknameExists
		}
	}
	return nil
}

func (r *UserRepository) FindByPk(userID int) (*model.User, error) {
	u, ok := r.users[userID]
	if !ok {
		return nil, store.ErrRecordNotFound
	}
	return u, nil
}

// FindByEmailOrNick extract user from test storage by email or nickname
func (r *UserRepository) FindByEmailOrNick(word string) (*model.User, error) {
	for _, u := range r.users {
		if u.Email == word || u.Nickname == word {
			return u, nil
		}
	}
	return nil, store.ErrRecordNotFound
}
