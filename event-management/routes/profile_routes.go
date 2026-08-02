package routes

import (
	"example.com/learning/event-management/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterProfileRoutes(router *gin.Engine, handler *handlers.ProfileHandler, authMiddleware gin.HandlerFunc) {
	profileRoutes := router.Group("/profile")

	profileRoutes.Use(authMiddleware)

	profileRoutes.GET("", handler.GetProfile)
	profileRoutes.PUT("", handler.UpdateProfile)
	profileRoutes.POST("/picture", handler.UploadProfilePicture)
	profileRoutes.DELETE("/picture", handler.DeleteProfilePicture)
}
