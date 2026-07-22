package routes

import (
	"example.com/learning/gin/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine, handler *handlers.UserHandler) {

	router.POST("/signup", handler.Signup)
	router.POST("/login", handler.Login)

}
