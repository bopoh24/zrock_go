package apiserver

import validation "github.com/go-ozzo/ozzo-validation"

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
