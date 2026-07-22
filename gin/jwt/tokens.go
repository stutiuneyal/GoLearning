package tokens

import (
	"errors"
	"fmt"
	"time"

	"example.com/learning/gin/constants"
	"example.com/learning/gin/models"
	"github.com/golang-jwt/jwt/v5"
)

// GenerateTokens -> generate a JWT for authenticated users
func GenerateTokens(user models.User) (string, error) {

	claims := jwt.MapClaims{
		"userId": user.Id,
		"email":  user.Email,
		"exp":    time.Now().Add(2 * time.Hour).Unix(),
		"iat":    time.Now().Unix(),
	}

	// Create the token using HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign this token using the secret key
	return token.SignedString([]byte(constants.SecretKey))

}

// VerifyToken verifies a JWT and returns the authenticated user's Id
func VerifyToken(tokenString string) (int, error) {

	if tokenString == "" {
		return 0, fmt.Errorf("token is required")
	}

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) { // This function basically returns the secret key used for signing the token

			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected signign algorithm: %s", token.Method.Alg())
			}

			return []byte(constants.SecretKey), nil
		},

		// explicit restriction -> that only HS256 is accepted by the parser
		jwt.WithValidMethods([]string{
			jwt.SigningMethodHS256.Alg(),
		}),
	)

	if err != nil {
		return 0, fmt.Errorf("parse token: %w", err)
	}

	if !token.Valid {
		return 0, errors.New("token is invalid")
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("token claims are invalid")
	}

	userIDValue, exists := mapClaims["userId"]
	if !exists {
		return 0, errors.New("userId claim is missing")
	}

	userIdFloat, ok := userIDValue.(float64)
	if !ok {
		return 0, errors.New("userId has invalid type")
	}

	userId := int(userIdFloat)

	if userId <= 0 {
		return 0, errors.New("userId claim is invalid")
	}

	return userId, nil

}
