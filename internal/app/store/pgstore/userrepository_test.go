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
	print(u)
	user, err = s.User().FindByPk(u.ID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, teardown := TestDB(t, databaseURL)
	defer teardown("users")
	s := New(db)
	email := "user@example.org"
	user, err := s.User().FindByEmail(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())
	assert.Nil(t, user)
	uFixture := model.UserFixture()
	uFixture.Email = email
	s.User().Create(uFixture)
	user, err = s.User().FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestUserRepository_FindByNick(t *testing.T) {
	db, teardown := TestDB(t, databaseURL)
	defer teardown("users")
	s := New(db)
	nick := "cool_hacker"
	user, err := s.User().FindByNick(nick)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())
	assert.Nil(t, user)
	uFixture := model.UserFixture()
	uFixture.Nick = nick
	s.User().Create(uFixture)
	user, err = s.User().FindByNick(nick)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}
