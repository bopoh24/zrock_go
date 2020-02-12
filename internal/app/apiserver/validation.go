package apiserver

import (
	"github.com/bopoh24/zrock_go/internal/app/model"
	validation "github.com/go-ozzo/ozzo-validation"
)

type errorResp struct {
	Error interface{}
}

type registerData struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Nickname  string `json:"nickname"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
}

type loginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (l *loginData) validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.Username, validation.Required),
		validation.Field(&l.Password, validation.Required),
	)
}

type loginResponseData struct {
	*model.User
	Token string `json:"token"`
}
