package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"example.com/learning/event-management/db/query"
	"example.com/learning/event-management/models"
	"golang.org/x/crypto/bcrypt"
)

var _ UserRepository = (*UserRepositoryImpl)(nil)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
)

type UserRepository interface {
	Save(ctx context.Context, user *models.User) error
	Login(ctx context.Context, user *models.User) error
	GetProfile(ctx context.Context, userID int) (models.Profile, error)
	UpdateProfile(ctx context.Context, userID int, name, bio *string) (models.Profile, error)
	SetProfilePicture(ctx context.Context, userID int, newObjectKey string) (oldObjectKey string, err error)
	RemoveProfilePicture(ctx context.Context, userID int) (oldObjectKey string, err error)
}

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepositoryImpl(db *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (u *UserRepositoryImpl) Save(ctx context.Context, user *models.User) error {

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := tx.QueryRow(query.SignupUserQuery, user.Email, user.Password, user.Name, user.Bio).Scan(&user.Id); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}

func (u *UserRepositoryImpl) Login(ctx context.Context, user *models.User) error {

	var savedUser models.User

	err := u.db.QueryRow(query.GetUserByEmail, user.Email).Scan(&savedUser.Id, &savedUser.Email, &savedUser.Password, &user.Name, &user.Bio, &user.ProfilePictureKey)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvalidCredentials
		}

		return fmt.Errorf("get user by email: %w", err)
	}

	// compare the hashed passwords
	if err := bcrypt.CompareHashAndPassword([]byte(savedUser.Password), []byte(user.Password)); err != nil {
		return err
	}

	*user = savedUser

	return nil

}

func (u *UserRepositoryImpl) GetProfile(ctx context.Context, userID int) (models.Profile, error) {
	var profile models.Profile

	err := u.db.QueryRowContext(ctx, query.GetProfileByUserIDQuery, userID).Scan(&profile.UserID, &profile.Email, &profile.Name, &profile.Bio, &profile.ProfilePictureKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Profile{}, ErrUserNotFound
		}

		return models.Profile{}, fmt.Errorf(
			"get profile: %w",
			err,
		)
	}

	return profile, nil
}

func (u *UserRepositoryImpl) UpdateProfile(ctx context.Context, userID int, name, bio *string) (models.Profile, error) {
	var profile models.Profile

	nameWasProvided := name != nil
	bioWasProvided := bio != nil

	nameValue := ""
	if name != nil {
		nameValue = *name
	}

	bioValue := ""
	if bio != nil {
		bioValue = *bio
	}

	err := u.db.QueryRowContext(ctx, query.UpdateProfileQuery, userID, nameWasProvided, nameValue, bioWasProvided, bioValue).Scan(&profile.UserID, &profile.Email, &profile.Name, &profile.Bio, &profile.ProfilePictureKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Profile{}, ErrUserNotFound
		}

		return models.Profile{}, fmt.Errorf(
			"update profile: %w",
			err,
		)
	}

	return profile, nil
}

func (u *UserRepositoryImpl) SetProfilePicture(ctx context.Context, userID int, newObjectKey string) (string, error) {

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf(
			"begin profile-picture transaction: %w",
			err,
		)
	}

	defer tx.Rollback()

	var oldObjectKey string

	err = tx.QueryRowContext(ctx, query.GetProfilePictureForUpdateQuery, userID).Scan(&oldObjectKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrUserNotFound
		}

		return "", fmt.Errorf(
			"get existing profile picture: %w",
			err,
		)
	}

	result, err := tx.ExecContext(ctx, query.UpdateProfilePictureQuery, userID, newObjectKey)
	if err != nil {
		return "", fmt.Errorf(
			"update profile picture: %w",
			err,
		)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf(
			"read profile-picture update result: %w",
			err,
		)
	}

	if affected == 0 {
		return "", ErrUserNotFound
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf(
			"commit profile-picture update: %w",
			err,
		)
	}

	return oldObjectKey, nil
}

func (u *UserRepositoryImpl) RemoveProfilePicture(ctx context.Context, userID int) (string, error) {

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf(
			"begin profile-picture removal transaction: %w",
			err,
		)
	}

	defer tx.Rollback()

	var oldObjectKey string

	err = tx.QueryRowContext(ctx, query.GetProfilePictureForUpdateQuery, userID).Scan(&oldObjectKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrUserNotFound
		}

		return "", fmt.Errorf(
			"get existing profile picture: %w",
			err,
		)
	}

	result, err := tx.ExecContext(ctx, query.RemoveProfilePictureQuery, userID)
	if err != nil {
		return "", fmt.Errorf(
			"remove profile picture: %w",
			err,
		)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf(
			"read profile-picture removal result: %w",
			err,
		)
	}

	if affected == 0 {
		return "", ErrUserNotFound
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf(
			"commit profile-picture removal: %w",
			err,
		)
	}

	return oldObjectKey, nil
}
