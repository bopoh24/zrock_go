package store

import "errors"

var (
	// ErrRecordNotFound ..
	ErrRecordNotFound = errors.New("record not found")
	// ErrEmailExists error if uer with this email already exists
	ErrEmailExists = errors.New("email is already taken")
	// ErrNicknameExists error if uer with this nickname already exists
	ErrNicknameExists = errors.New("nickname is already taken")
)
