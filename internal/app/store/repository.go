package store

import "github.com/bopoh24/zrock_go/internal/app/model"

// UserRepository interface
type UserRepository interface {
	Create(*model.User) error
	FindByPk(int) (*model.User, error)
	FindByEmailOrNick(string) (*model.User, error)
	VerifyEmail(string, string) error
}
