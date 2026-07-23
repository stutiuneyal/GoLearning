package main

import (
	"fmt"
	"os"

	"example.com/learning/gin/config"
	database "example.com/learning/gin/db"
	"example.com/learning/gin/handlers"
	"example.com/learning/gin/repository"
	"example.com/learning/gin/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	// connect to the database
	db, err := config.ConnectToDatabase()
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// create tables
	if err := database.CreateTables(db); err != nil {
		fmt.Printf("Error creating tables: %v\n", err)
		os.Exit(1)
	}

	// register the repositories
	userRepo := repository.NewUserRepositoryImpl(db)
	eventRepo := repository.NewEventRepositoryImpl(db)

	// register the handlers
	userHandler := handlers.NewUserHandler(userRepo)
	eventHandler := handlers.NewEventHandler(eventRepo)

	// setting up the gin server
	server := gin.Default()

	// register the routes
	routes.RegisterUserRoutes(server, userHandler)
	routes.RegisterEventRoutes(server, eventHandler)

	if err := server.Run(":8080"); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(2)
	}

}
