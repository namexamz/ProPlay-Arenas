package config

import (
	"log"
	"os"
)

// EnsureLogDir создает директорию для логов, если её нет.
// Читает путь из переменной окружения LOG_DIR, по умолчанию "logs".
func EnsureLogDir() string {
	dir := os.Getenv("LOG_DIR")
	if dir == "" {
		dir = "logs"
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		log.Printf("failed to create log dir %s: %v", dir, err)
	}
	return dir
}
