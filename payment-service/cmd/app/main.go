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

	// Инициализация логгера
	config.InitLogger()

	// Подключение к БД
	db := config.ConnectDB()

	// Миграция моделей
	if err := db.AutoMigrate(
		&models.Payment{},
		&models.Refund{},
	); err != nil {
		config.Error("ошибка миграции базы данных", err)
	}

	config.Info("✓ миграция БД успешна")
	config.Info("✓ сервис платежей готов")
}
