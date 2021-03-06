package apiserver

import "errors"

type usernameError struct {
	Username string `json:"username"`
}

func (e *usernameError) Error() string {
	return e.Username
}

var (
	ErrUsernameOrPassword = &usernameError{"incorrect username or password"}
	ErrEmailNotVerified   = &usernameError{"email not verified"}
	JSONDecodeError       = errors.New("incorrect JSON recieved")
)
