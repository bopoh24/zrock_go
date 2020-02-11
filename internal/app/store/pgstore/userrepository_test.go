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

	// creating user with exists тшслтфьу
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
