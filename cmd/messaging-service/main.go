package cmd

import (
	"context"
	"log"
	"messaging-service/internal/repository"
	"messaging-service/internal/service"
	"messaging-service/pkg/database"
	"messaging-service/pkg/logger"
	"messaging-service/utils"
)

func Execute() {
	config, err := utils.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	logger.InitLogger(config.LogLevel)
	logger.Info("Logger initialized")
	db, err := database.NewMongoDB(config.DatabaseURI, config.DatabaseName)
	if err != nil {
		logger.Fatal("Failed to connect to database: %v", err)
	}
	defer db.Close()
	logger.Info("Connected to database")
	logger.Info("Initializing services")
	usersCollection := db.Collection(config.UsersCollection)

	userRepo := repository.NewUserRepository(usersCollection)
	userService := service.NewUserService(userRepo)
	logger.Info("Services initialized")
	tempCtx := context.Background()
	userService.GetUserByID(tempCtx, "123")
	userService.GetUserByEmail(tempCtx, "123")
	logger.Info("User fetched")
}
