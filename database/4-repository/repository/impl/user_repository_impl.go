package impl

import (
	"context"
	"database/sql"

	"example.com/learning/database/4-repository/models"
	"example.com/learning/database/4-repository/repository"
	"golang.org/x/crypto/bcrypt"
)

var _ repository.UserRepository = (*UserRepositoryImpl)(nil)

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepositoryImpl(db *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (u *UserRepositoryImpl) CreateUser(name, email, password, avatar string) (int64, error) {
	ctx := context.Background()

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return 0, err
	}

	userInsertQuery := `
		INSERT INTO users (name, email, hashed_password)
		VALUES ($1,$2,$3)
		RETURNING id;
	`

	user := models.User{
		Name:           name,
		Email:          email,
		HashedPassword: string(hashedPassword),
	}

	if err := u.db.QueryRow(userInsertQuery, user.Name, user.Email, user.HashedPassword).Scan(&user.Id); err != nil {
		return 0, err
	}

	// user is created successfully -> proceed and create user profile

	userProfileInsertQuery := `
		INSERT INTO profiles (user_id,avatar)
		VALUES ($1,$2)
	`

	if _, err := u.db.Exec(userProfileInsertQuery, user.Id, avatar); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return int64(user.Id), nil
}

func (u *UserRepositoryImpl) GetUserById(id int64) (models.User, error) {
	getUserByIdQuery := `
		SELECT u.id, u.name, u.email, u.hashed_password, p.avatar
		FROM users u INNER JOIN profiles p
		ON u.id = p.user_id
		WHERE u.id=$1;
	`

	var user models.User

	if err := u.db.QueryRow(getUserByIdQuery, id).Scan(&user.Id, &user.Name, &user.Email, &user.HashedPassword, &user.Profile.Avatar); err != nil {
		return models.User{}, err
	}

	return user, nil
}
