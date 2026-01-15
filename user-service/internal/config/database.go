package config

import (
	"fmt"
	"log"
	"os"
	"user-service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	if dbHost == "" || dbUser == "" || dbName == "" || dbPort == "" {
		log.Fatal("One or more required environment variables are missing: DB_HOST, DB_USER, DB_NAME, DB_PORT")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s dbname=%s port=%s ",
		dbHost,
		dbUser,
		dbName,
		dbPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	// Автоматическая миграция моделей
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("failed to auto-migrate models: ", err)
	}

	DB = db
}
