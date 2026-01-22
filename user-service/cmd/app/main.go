package main

import (
	"os"

	"log/slog"
	"user-service/internal/config"
	"user-service/internal/repository"
	service "user-service/internal/services"
	"user-service/internal/transport"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// env
	if os.Getenv("ENV") != "production" {
		_ = godotenv.Load()
	}

	// logger
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	).With(slog.String("service", "user-service"))

	db := config.ConnectDatabase()

	// repos & services
	userRepo := repository.NewUserRepository(logger, db)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		logger.Error("JWT_SECRET is not set")
		os.Exit(1)
	}

	authService := service.NewAuthService(logger, jwtSecret, userRepo)
	userService := service.NewUserService(logger, userRepo)

	// handlers
	authHandler := transport.NewAuthHandler(logger, authService)
	userHandler := transport.NewUserHandler(userService, logger)

	// gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.SetTrustedProxies(nil)

	api := r.Group("/api")

	auth := api.Group("/auth")
	authHandler.RegisterRoutes(auth)

	protected := api.Group("/")
	// protected.Use(transport.AuthMiddleware(jwtSecret))
	userHandler.RegisterRoutes(protected)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("server started", slog.String("port", port))
	

	if err := r.Run(":" + port); err != nil {
		logger.Error("server failed", slog.Any("error", err))
	}


}
