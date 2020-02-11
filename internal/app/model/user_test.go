package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_BeforeCreate(t *testing.T) {
	u := UserFixture()
	assert.NoError(t, u.BeforeCreate())
	assert.NotEmpty(t, u.EncryptedPassword)
}

func TestUser_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		u       func() *User
		isValid bool
	}{
		{
			name: "valid",
			u: func() *User {
				return UserFixture()
			},
			isValid: true,
		},
		{
			name: "empty email",
			u: func() *User {
				u := UserFixture()
				u.Email = ""
				return u
			},
			isValid: false,
		},
		{
			name: "invalid email",
			u: func() *User {
				u := UserFixture()
				u.Email = "coolMyMail"
				return u
			},
			isValid: false,
		},
		{
			name: "empty password",
			u: func() *User {
				u := UserFixture()
				u.Password = ""
				return u
			},
			isValid: false,
		},
		{
			name: "short password",
			u: func() *User {
				u := UserFixture()
				u.Password = "short"
				return u
			},
			isValid: false,
		},
		{
			name: "with encrypted password",
			u: func() *User {
				u := UserFixture()
				u.Password = ""
				u.EncryptedPassword = "someEncryptedPassordData"
				return u
			},
			isValid: true,
		},
		{
			name: "no nick",
			u: func() *User {
				u := UserFixture()
				u.Nickname = ""
				return u
			},
			isValid: false,
		},
		{
			name: "nick with unallowed symbols",
			u: func() *User {
				u := UserFixture()
				u.Nickname = "!!!cool!!!"
				return u
			},
			isValid: false,
		},
		{
			name: "no name",
			u: func() *User {
				u := UserFixture()
				u.FirstName = ""
				return u
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.u().Validate())
			} else {
				assert.Error(t, tc.u().Validate())
			}
		})
	}
}
