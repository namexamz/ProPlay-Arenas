package main

import (
	"log"
	"payment-service/internal/config"
	"payment-service/internal/models"

	"github.com/joho/godotenv"
)

func main() {
	// Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		log.Println("не удалось загрузить переменные окружения", err)
	}
	// Подключение к БД
	db := config.ConnectDB()

	// Миграция моделей
	if err := db.AutoMigrate(
		&models.Payment{},
		&models.Refund{},
	); err != nil {
		log.Fatalf("ошибка миграции базы данных: %v", err)
	}

	log.Println("✓ миграция БД успешна")
	log.Println("✓ сервис платежей готов")
}
