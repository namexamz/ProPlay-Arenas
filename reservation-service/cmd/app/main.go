package main

import (
	"log"
	"os"
	"reservation/internal/config"
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
		log.Fatal("Warning: .env file not found, using system environment variables", err)
	}

	db := config.SetUpDatabaseConnection()

	if err := db.AutoMigrate(&models.ReservationDetails{}); err != nil {
		log.Fatal("Ошибка миграции базы данных:", err)
	}

	bookingRepo := repository.NewBookingRepo(db)
	bookingServ := service.NewBookingServ(bookingRepo)

	// Получение JWT секретного ключа
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET не задан в переменных окружения")
	}

	r := gin.Default()

	bookingHandler := transport.NewBookingHandler(r, bookingServ)
	bookingHandler.Register(r, jwtSecret)

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("Сервер запущен на порту %s", port)
	r.Run(":" + port)
}
