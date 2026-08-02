package handlers

import (
	"net/http"

	"example.com/learning/event-management/dto"
	tokens "example.com/learning/event-management/jwt"
	"example.com/learning/event-management/models"
	"example.com/learning/event-management/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	UserRepository repository.UserRepository
	jwtSecret      string
}

func NewUserHandler(userRepo repository.UserRepository, jwtSecret string) *UserHandler {
	return &UserHandler{
		UserRepository: userRepo,
		jwtSecret:      jwtSecret,
	}
}

func (u *UserHandler) Signup(c *gin.Context) {

	var request dto.SignupRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	// hash the password
	hashed, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	user := models.User{
		Email:    request.Email,
		Password: string(hashed),
		Name:     "",
		Bio:      "",
	}

	// save the user
	if err := u.UserRepository.Save(c.Request.Context(), &user); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	userResponse := dto.UserResponse{
		Id:                user.Id,
		Email:             user.Email,
		Name:              user.Name,
		Bio:               user.Bio,
		ProfilePictureURL: nil,
	}

	c.JSON(http.StatusCreated, gin.H{"user": userResponse})

}

func (u *UserHandler) Login(c *gin.Context) {

	var request dto.LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	user := models.User{
		Email:    request.Email,
		Password: request.Password,
	}

	if err := u.UserRepository.Login(c.Request.Context(), &user); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	token, err := tokens.GenerateTokens(user, u.jwtSecret)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "failed to generate authentication token",
			},
		)
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}
