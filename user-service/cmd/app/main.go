package main

import (
	"log"
	"os"

	"log/slog"
	"user-service/internal/config"
	"user-service/internal/middleware"
	"user-service/internal/repository"
	service "user-service/internal/services"
	"user-service/internal/transport"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Загрузка переменных окружения из .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла: ", err)
	}

	// Инициализация базы данных
	config.ConnectDatabase()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	db := config.DB

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(logger, db)

	// Получение секретного ключа из переменных окружения
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET не задан в переменных окружения")
	}

	// Инициализация сервисов
	authService := service.NewAuthService(jwtSecret, userRepo)
	userService := service.NewUserService(logger, userRepo)

	// Инициализация обработчиков
	authHandler := transport.NewAuthHandler(authService)
	userHandler := transport.NewUserHandler(userService, logger)

	// Настройка маршрутизатора
	r := gin.Default()

	// Маршруты без авторизации
	authGroup := r.Group("/auth")
	authHandler.RegisterRoutes(authGroup)

	// Маршруты с авторизацией
	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware(authService))
	userHandler.RegisterRoutes(authorized)

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Сервер запущен на порту %s", port)
	r.Run(":" + port)
}
