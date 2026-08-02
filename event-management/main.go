package main

import (
	"context"
	"fmt"
	"os"

	"example.com/learning/event-management/config"
	database "example.com/learning/event-management/db"
	"example.com/learning/event-management/handlers"
	"example.com/learning/event-management/middleware"
	"example.com/learning/event-management/repository"
	"example.com/learning/event-management/routes"
	"example.com/learning/event-management/storage"
	"github.com/gin-gonic/gin"
)

func main() {

	// load the config
	appConfig, err := config.Load()
	if err != nil {
		fmt.Printf("Configuration Error: %v\n", err)
		os.Exit(1)
	}

	// connect to the database
	db, err := config.ConnectToDatabase(appConfig)
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

	// Configure local object storage for profile pictures.
	objectStorage, localStorage, err := buildObjectStorage(context.Background(), appConfig)
	if err != nil {
		fmt.Printf("Storage error: %w\n", err)
		os.Exit(1)
	}

	// register the handlers
	userHandler := handlers.NewUserHandler(userRepo, appConfig.JWTSecret)
	eventHandler := handlers.NewEventHandler(eventRepo)
	profileHandler := handlers.NewProfileHandler(userRepo, objectStorage)

	// register middleware
	authMiddleware := middleware.Authenticate(appConfig.JWTSecret)

	// setting up the gin server
	server := gin.Default()

	// Limit the amount of multipart data retained in memory.
	//
	// Uploaded profile pictures are limited separately to 5 MB by the handler.
	server.MaxMultipartMemory = 8 << 20

	// Expose locally uploaded files over HTTP.
	if localStorage != nil {
		server.Static(appConfig.UploadURLBase, localStorage.RootDir())
	}

	// register the routes
	routes.RegisterHealthRoutes(server)
	routes.RegisterUserRoutes(server, userHandler)
	routes.RegisterEventRoutes(server, eventHandler, authMiddleware)
	routes.RegisterProfileRoutes(server, profileHandler, authMiddleware)

	listenAddress := ":" + appConfig.ServerPort

	fmt.Printf("Event Management API listening on %s\n", listenAddress)

	if err := server.Run(listenAddress); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(2)
	}

}

func buildObjectStorage(ctx context.Context, appConfig config.Config) (storage.Store, *storage.LocalStore, error) {
	switch appConfig.StorageProvider {
	case "local":
		localStorage, err := storage.NewLocalStore(appConfig.UploadDirectory, appConfig.UploadURLBase)
		if err != nil {
			return nil, nil, err
		}
		return localStorage, localStorage, nil
	case "s3":
		s3Storage, err := storage.NewS3Store(ctx, appConfig.AWSRegion, appConfig.S3Bucket)
		if err != nil {
			return nil, nil, err
		}
		return s3Storage, nil, nil
	default:
		return nil, nil, fmt.Errorf(
			"unsupported storage provider: %s",
			appConfig.StorageProvider,
		)
	}
}
