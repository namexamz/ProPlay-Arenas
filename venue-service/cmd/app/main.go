package main

import (
	"fmt"
	"log"
	"net/http"

	"venue-service/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Ошибка загрузки переменных окружения: ", err)
	}

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("ConnectDB: %v", err)
	}
	_ = db // DB будет использоваться в handlers/repository/services

	r := gin.Default()

	// Отключаем доверие прокси для локальной разработки
	// Для production используйте: r.SetTrustedProxies([]string{"127.0.0.1"})
	// TODO: не забыть включить доверие прокси для production
	r.SetTrustedProxies(nil)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})

	if err := r.Run(fmt.Sprintf(":%s", config.GetEnv("PORT", "8080"))); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
