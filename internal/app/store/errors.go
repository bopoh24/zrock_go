package store

import "errors"

type emailError struct {
	Email string `json:"email"`
}

func (e *emailError) Error() string {
	return e.Email
}

type nicknameError struct {
	Nickname string `json:"nickname"`
}

func (e *nicknameError) Error() string {
	return e.Nickname
}

var (
	// ErrRecordNotFound if no record found
	ErrRecordNotFound = errors.New("record not found")
	// ErrEmailExists error if user with this email already exists
	ErrEmailExists = &emailError{"email is already taken"}
	// ErrNicknameExists error if uer with this nickname already exists
	ErrNicknameExists = &nicknameError{"nickname is already taken"}
)
