package main

import (
	"log"
	"os"
	"reservation/internal/config"
	"reservation/internal/kafka"
	"reservation/internal/models"
	"reservation/internal/repository"
	"reservation/internal/service"
	"reservation/internal/transport"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(".env file not found, using system environment variables", err)
	}

	db := config.SetUpDatabaseConnection()

	if err := db.AutoMigrate(&models.ReservationDetails{}); err != nil {
		log.Fatal("Ошибка миграции базы данных:", err)
	}

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		log.Fatal("KAFKA_BROKERS не задан в переменных окружения")
	}

	producer := kafka.NewProducer([]string{kafkaBrokers})
	defer func() {
		if err := producer.Close(); err != nil {
			log.Printf("Ошибка закрытия Kafka продюсера: %v", err)
		}
	}()

	bookingRepo := repository.NewBookingRepo(db)
	bookingServ := service.NewBookingServ(bookingRepo, producer)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET не задан в переменных окружения")
	}

	r := gin.Default()

	transport.RegisterRoutes(r, bookingServ, jwtSecret)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("Сервер запущен на порту %s", port)
	r.Run(":" + port)
}
