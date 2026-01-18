package config

import (
	"log"
	"os"
)

var Logger *log.Logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

func InitLogger() *log.Logger {
	// Определяем выходной файл для логов
	logFile := os.Getenv("LOG_FILE")
	if logFile == "" {
		logFile = "payment-service.log"
	}

	// Открываем или создаем файл логов
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("ошибка открытия файла логов: %v", err)
	}

	// Создаем логгер с указанием файла и флагов
	Logger = log.New(
		file,
		"",
		log.LstdFlags|log.Lshortfile,
	)

	return Logger
}

// Info логирует информационные сообщения
func Info(msg string) {
	if Logger == nil {
		return
	}
	Logger.Printf("[INFO] %s", msg)
}

// Error логирует ошибки
func Error(msg string, err error) {
	if Logger == nil {
		return
	}
	if err != nil {
		Logger.Printf("[ERROR] %s: %v", msg, err)
	} else {
		Logger.Printf("[ERROR] %s", msg)
	}
}

// Warn логирует предупреждения
func Warn(msg string) {
	if Logger == nil {
		return
	}
	Logger.Printf("[WARN] %s", msg)
}

// Debug логирует отладочные сообщения
func Debug(msg string) {
	if Logger == nil {
		return
	}
	if os.Getenv("DEBUG") == "true" {
		Logger.Printf("[DEBUG] %s", msg)
	}
}
