package model

import (
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// User model
type User struct {
	ID                int        `json:"id"`
	Email             string     `json:"email"`
	Password          string     `json:"password,omitempty"`
	EncryptedPassword string     `json:"-"`
	Nickname          string     `json:"nickname"`
	FirstName         string     `json:"first_name"`
	LastName          string     `json:"last_name,omitempty"`
	Avatar            string     `json:"avatar,omitempty"`
	LastLogin         *time.Time `json:"last_login,omitempty"`
	Created           *time.Time `json:"created,omitempty"`
}

// BeforeCreate before user model create
func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := encryptString(u.Password)
		if err != nil {
			return err
		}
		u.EncryptedPassword = enc
	}
	return nil
}

// Validate user model
func (u *User) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Nickname, validation.Required, validation.Length(3, 32),
			validation.Match(regexp.MustCompile("^[a-zA-Z0-9_]+$")).Error("only latin letters, numbers and underscores are allowed")),
		validation.Field(&u.FirstName, validation.Required, validation.Length(2, 32)),
		validation.Field(&u.Password, validation.By(requiredIf(u.EncryptedPassword == "")),
			validation.Length(6, 40),
			validation.Match(regexp.MustCompile(`^[^\s\r\n]+$`)).Error("space is not allowed")))
}

// ComparePassword checks if user password is correct
func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

// Sanitize clears all restricted data
func (u *User) Sanitize() {
	u.Password = ""
}

func encryptString(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
