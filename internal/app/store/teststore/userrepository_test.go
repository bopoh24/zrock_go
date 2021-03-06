package teststore

import (
	"testing"

	"github.com/bopoh24/zrock_go/internal/app/model"
	"github.com/bopoh24/zrock_go/internal/app/store"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	s := New()
	user := model.UserFixture()
	assert.NoError(t, s.User().Create(user))

	// creating user with exists email
	assert.Equal(t, store.ErrEmailExists, s.User().Create(user))

	// creating user with exists nickname
	user = model.UserFixture()
	user.Email = "new@email.org"
	assert.Equal(t, store.ErrNicknameExists, s.User().Create(user))
}

func TestUserRepository_FindByPk(t *testing.T) {
	s := New()
	user, err := s.User().FindByPk(1)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())
	assert.Nil(t, user)
	u := model.UserFixture()
	s.User().Create(u)
	user, err = s.User().FindByPk(u.ID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestUserRepository_FindByEmailOrNick(t *testing.T) {
	email := "user@example.org"
	s := New()
	user, err := s.User().FindByEmailOrNick(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())
	assert.Nil(t, user)
	u := model.UserFixture()
	u.Email = email
	s.User().Create(u)
	user, err = s.User().FindByEmailOrNick(email)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	u.Email = "user@example2.org"
	u.Nickname = "hacker"
	s.User().Create(u)
	user, err = s.User().FindByEmailOrNick("some_nick")
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())
	assert.Nil(t, user)
	user, err = s.User().FindByEmailOrNick(u.Nickname)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestVerifyEmail(t *testing.T) {
	uFixture := model.UserFixture()
	s := New()
	s.User().Create(uFixture)
	testCases := []struct {
		name             string
		email            string
		verificationCode string
		errorExpected    bool
	}{
		{
			name:             "incorrect email and code",
			email:            "",
			verificationCode: "",
			errorExpected:    true,
		},
		{
			name:             "incorrect email",
			email:            "",
			verificationCode: uFixture.EmailVerificationCode,
			errorExpected:    true,
		},
		{
			name:             "incorrect code",
			email:            uFixture.Email,
			verificationCode: "",
			errorExpected:    true,
		},
		{
			name:             "correct email and code",
			email:            uFixture.Email,
			verificationCode: uFixture.EmailVerificationCode,
			errorExpected:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := s.User().VerifyEmail(tc.email, tc.verificationCode)
			if tc.errorExpected {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
