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

func TestUserRepository_FindByEmail(t *testing.T) {
	email := "user@example.org"
	s := New()
	user, err := s.User().FindByEmail(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())
	assert.Nil(t, user)
	s.User().Create(model.UserFixture())
	user, err = s.User().FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestUserRepository_FindByNick(t *testing.T) {
	nick := "hacker"
	s := New()
	user, err := s.User().FindByNick(nick)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())
	assert.Nil(t, user)
	userModel := model.UserFixture()
	userModel.Nick = nick
	s.User().Create(userModel)
	user, err = s.User().FindByNick(nick)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}
