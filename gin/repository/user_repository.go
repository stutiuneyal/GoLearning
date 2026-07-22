package repository

import (
	"context"
	"database/sql"
	"fmt"

	"example.com/learning/gin/db/query"
	"example.com/learning/gin/models"
	"golang.org/x/crypto/bcrypt"
)

var _ UserRepository = (*UserRepositoryImpl)(nil)

type UserRepository interface {
	Save(user *models.User) error
	Login(user *models.User) error
}

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepositoryImpl(db *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (u *UserRepositoryImpl) Save(user *models.User) error {

	ctx := context.Background()

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := u.db.QueryRow(query.SignupUserQuery, user.Email, user.Password).Scan(&user.Id); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}

func (u *UserRepositoryImpl) Login(user *models.User) error {

	var savedUser models.User

	if err := u.db.QueryRow(query.GetUserByEmail, user.Email).Scan(&savedUser.Id, &savedUser.Email, &savedUser.Password); err != nil {
		return err
	}

	if savedUser.Email != user.Email || savedUser.Password == "" {
		return fmt.Errorf("invalid user")
	}

	// compare the hashed passwords
	if err := bcrypt.CompareHashAndPassword([]byte(savedUser.Password), []byte(user.Password)); err != nil {
		return err
	}

	user.Id = savedUser.Id
	user.Password = savedUser.Password

	return nil

}
