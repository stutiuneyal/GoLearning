package middleware

import (
	"net/http"
	"strings"

	tokens "example.com/learning/gin/jwt"
	"github.com/gin-gonic/gin"
)

func Authenticate(c *gin.Context) {

	// read the header from the request
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"error": "authorization header is required"},
		)
		return
	}

	// Token -> "Bearer eyguyhb....."
	parts := strings.Fields(authHeader)

	if len(parts) != 2 {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"error": "authorization header is not in expected format"},
		)
		return
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"error": "authorization token is not bearer"},
		)
		return
	}

	tokenString := parts[1]

	// verify the token and retreve the authenticated user Id
	userId, err := tokens.VerifyToken(tokenString)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"error": "invalid or expired token"},
		)
		return
	}

	// store the authenticated userId in the gin's context
	c.Set("userId", userId)

	c.Next()

}
