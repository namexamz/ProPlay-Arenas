package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"payment-service/internal/config"
	"payment-service/internal/models"
	"payment-service/internal/repository"
	"payment-service/internal/services"
	"payment-service/internal/transport"
	kafkaconsumer "payment-service/internal/transport/kafka"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		slog.Warn("не удалось загрузить .env", "error", err)
	}

	db := config.ConnectDB()

	if err := db.AutoMigrate(
		&models.Payment{},
		&models.Refund{},
	); err != nil {
		slog.Error("ошибка миграции схемы", "error", err)
		os.Exit(1)
	}

	paymentRepo := repository.NewPaymentRepository(db)
	refundRepo := repository.NewRefundRepository(db)
	paymentService := services.NewPaymentService(paymentRepo)
	refundService := services.NewRefundService(refundRepo, paymentRepo, db)
	transportHandler := transport.NewPaymentHandler(paymentService, refundService, logger)
	consumer := kafkaconsumer.NewConsumerFromEnv(paymentService, refundService, logger)
	consumer.Start(context.Background())

	r := gin.Default()
	api := r.Group("/")
	transportHandler.RegisterRoutes(api)

	port := config.GetEnv("PORT", "8084")
	slog.Info("HTTP сервер запущен", "port", port)
	if err := r.Run(":" + port); err != nil {
		slog.Error("не удалось запустить HTTP сервер", "error", err)
		os.Exit(1)
	}
}
