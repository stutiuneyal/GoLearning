package repository

import (
	"example.com/learning/database/4-repository/models"
)

type UserRepository interface {
	CreateUser(name, email, password, avatar string) (int64, error)
	GetUserById(id int64) (models.User, error)
}
