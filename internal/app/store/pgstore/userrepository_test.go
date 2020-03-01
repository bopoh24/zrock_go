package pgstore

import (
	"testing"

	"github.com/bopoh24/zrock_go/internal/app/model"
	"github.com/bopoh24/zrock_go/internal/app/store"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := TestDB(t, databaseURL)
	defer teardown("users")
	s := New(db)
	user := model.UserFixture()
	assert.NoError(t, s.User().Create(user))

	// creating user with exists email
	user = model.UserFixture()
	assert.Equal(t, store.ErrEmailExists, s.User().Create(user))

	// creating user with exists nickname
	user.Email = "new@email.org"
	assert.Equal(t, store.ErrNicknameExists, s.User().Create(user))
}

func TestUserRepository_FindByPk(t *testing.T) {
	db, teardown := TestDB(t, databaseURL)
	defer teardown("users")
	s := New(db)
	id := 0
	user, err := s.User().FindByPk(id)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())
	assert.Nil(t, user)
	u := model.UserFixture()
	s.User().Create(u)
	user, err = s.User().FindByPk(u.ID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestUserRepository_FindByEmailOrNick(t *testing.T) {
	db, teardown := TestDB(t, databaseURL)
	defer teardown("users")
	s := New(db)
	email := "user@example.org"
	user, err := s.User().FindByEmailOrNick(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())
	assert.Nil(t, user)
	uFixture := model.UserFixture()
	uFixture.Email = email
	s.User().Create(uFixture)
	user, err = s.User().FindByEmailOrNick(email)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	uFixture.Email = "user@example2.org"
	uFixture.Nickname = "hacker"
	s.User().Create(uFixture)
	user, err = s.User().FindByEmailOrNick("some_nick")
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())
	assert.Nil(t, user)
	user, err = s.User().FindByEmailOrNick(uFixture.Nickname)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestVerifyEmail(t *testing.T) {
	db, teardown := TestDB(t, databaseURL)
	defer teardown("users")
	uFixture := model.UserFixture()
	s := New(db)
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
