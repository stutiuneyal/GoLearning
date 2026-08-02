package routes

import (
	"example.com/learning/event-management/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterEventRoutes(router *gin.Engine, eventHandler *handlers.EventHandler, authMiddleware gin.HandlerFunc) {

	eventGroup := router.Group("/events")

	eventGroup.Use(authMiddleware) // middleware registration

	eventGroup.GET("", eventHandler.GetEvents)
	eventGroup.POST("", eventHandler.CreateEvent)
	eventGroup.GET("/:id", eventHandler.GetEventById)
	eventGroup.PUT("/:id", eventHandler.Update)
	eventGroup.DELETE("/:id", eventHandler.Delete)
}
