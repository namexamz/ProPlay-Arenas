package main

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"payment-service/internal/config"
	"payment-service/internal/models"
	"payment-service/internal/repository"
	"payment-service/internal/services"
	"payment-service/internal/transport"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		slog.Warn("could not load .env", "error", err)
	}

	db := config.ConnectDB()

	if err := db.AutoMigrate(
		&models.Payment{},
		&models.Refund{},
	); err != nil {
		slog.Error("failed to auto-migrate schema", "error", err)
		os.Exit(1)
	}

	paymentRepo := repository.NewPaymentRepository(db)
	refundRepo := repository.NewRefundRepository(db)
	paymentService := services.NewPaymentService(paymentRepo)
	refundService := services.NewRefundService(refundRepo, paymentRepo)
	transportHandler := transport.NewPaymentHandler(paymentService, refundService, logger)

	r := gin.Default()
	api := r.Group("/")
	transportHandler.RegisterRoutes(api)

	port := config.GetEnv("PORT", "8080")
	slog.Info("HTTP server listening", "port", port)
	if err := r.Run(":" + port); err != nil {
		slog.Error("failed to start HTTP server", "error", err)
		os.Exit(1)
	}
}
