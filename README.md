# Play-Arenas

Play-Arenas — это платформа онлайн-бронирования спортивных площадок и арен. Система позволяет владельцам площадок сдавать их в аренду, а пользователям находить и бронировать необходимую спортивную площадку.

Технологии проекта

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-008ECF?style=for-the-badge&logo=go&logoColor=white)
![GORM](https://img.shields.io/badge/GORM-7A3E9D?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)
![Kafka](https://img.shields.io/badge/Kafka-231F20?style=for-the-badge&logo=apache-kafka&logoColor=white)
![Zookeeper](https://img.shields.io/badge/Zookeeper-FF6600?style=for-the-badge&logo=apache&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-000000?style=for-the-badge&logo=json-web-tokens&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![Docker Compose](https://img.shields.io/badge/Docker_Compose-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![API Gateway](https://img.shields.io/badge/API_Gateway-0F172A?style=for-the-badge)
![Microservices](https://img.shields.io/badge/Microservices-2563EB?style=for-the-badge)

Установка зависимостей
go mod download

Запуск проекта
docker-compose up --build

Запуск приложения
# Обычный запуск
docker-compose up --build

# Запуск в фоновом режиме
docker-compose up -d --build

# Запуск только Gateway локально
go run ./gateway/cmd/app/main.go

Заполнение базы данных тестовыми данными
go run ./venue-service/cmd/seed/main.go

Доступные команды
make run      # Запуск reservation-service (из каталога reservation-service)
make dev      # Локальный запуск с hot-reload (при наличии конфигурации air)
make lint     # Проверка кода линтером (reservation-service)
make fmt      # Форматирование кода (reservation-service)
make vet      # Статический анализ кода (через go vet, при необходимости)
make test     # Запуск тестов (если добавлены)

Команда проекта
- [tsuruevimran17](https://github.com/tsuruevimran17)
- [Idigov](https://github.com/Idigov) (Baisangur)
- [namexamz](https://github.com/namexamz) (Хамзат)
- [DjMariarty](https://github.com/DjMariarty)
