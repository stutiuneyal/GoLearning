package routes

import (
	"example.com/learning/gin/handlers"
	"example.com/learning/gin/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterEventRoutes(router *gin.Engine, eventHandler *handlers.EventHandler) {

	eventGroup := router.Group("/events")

	eventGroup.Use(middleware.Authenticate) // middleware registration

	eventGroup.GET("", eventHandler.GetEvents)
	eventGroup.POST("", eventHandler.CreateEvent)
	eventGroup.GET("/:id", eventHandler.GetEventById)
	eventGroup.PUT("/:id", eventHandler.Update)
	eventGroup.DELETE("/:id", eventHandler.Delete)
}
