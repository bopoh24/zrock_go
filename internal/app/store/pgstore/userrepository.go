package pgstore

import (
	"database/sql"

	"github.com/bopoh24/zrock_go/internal/app/model"
	"github.com/bopoh24/zrock_go/internal/app/store"
)

// UserRepository repository for users
type UserRepository struct {
	store *Store
}

// Create add new user to database
func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}
	if err := r.checkEmailAndNickFree(u); err != nil {
		return err
	}
	if err := u.BeforeCreate(); err != nil {
		return err
	}
	if err := r.store.db.QueryRow(
		"INSERT INTO users (email, nickname, first_name, last_name, enpass, email_verification_code) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		u.Email, u.Nickname, u.FirstName, u.LastName, u.EncryptedPassword, u.EmailVerificationCode,
	).Scan(&u.ID); err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) checkEmailAndNickFree(u *model.User) error {

	if _, err := r.FindByEmailOrNick(u.Email); err == nil {
		return store.ErrEmailExists
	}
	if _, err := r.FindByEmailOrNick(u.Nickname); err == nil {
		return store.ErrNicknameExists
	}
	return nil
}

// FindByPk find user by ID
func (r *UserRepository) FindByPk(userID int) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, nickname, first_name, last_name, avatar, last_login, enpass from users WHERE id = $1",
		userID,
	).Scan(&u.ID, &u.Email, &u.Nickname, &u.FirstName, &u.LastName, &u.Avatar, &u.LastLogin, &u.EncryptedPassword); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}

// FindByEmailOrNick extract user from database by email or nickname
func (r *UserRepository) FindByEmailOrNick(word string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, email_verified, nickname, first_name, last_name, avatar, last_login, enpass from users WHERE email = $1 or nickname = $1",
		word,
	).Scan(&u.ID, &u.Email, &u.EmailVerified, &u.Nickname, &u.FirstName, &u.LastName, &u.Avatar, &u.LastLogin, &u.EncryptedPassword); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}

// VerifyEmail makes email verification
func (r *UserRepository) VerifyEmail(email string, verificationCode string) error {
	res, err := r.store.db.Exec("UPDATE users set email_verified=true where email=$1 and email_verification_code=$2", email, verificationCode)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return store.ErrRecordNotFound
	}
	return nil
}
