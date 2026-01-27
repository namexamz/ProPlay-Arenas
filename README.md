# Play-Arenas

## О проекте

Play-Arenas — это платформа онлайн-бронирования спортивных площадок и арен.  
Система позволяет владельцам площадок сдавать их во временное пользование, а пользователям находить и бронировать подходящие площадки по району, типу и времени.

## Технологии проекта

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

## Установка и запуск

### Установка зависимостей

В каждом сервисе:

```bash
go mod download
```

### Запуск всего проекта через Docker

```bash
docker-compose up --build
```

Запуск в фоновом режиме:

```bash
docker-compose up -d --build
```

### Локальный запуск Gateway

```bash
go run ./gateway/cmd/app/main.go
```

## Сидирование данных

Заполнение базы данных тестовыми площадками:

```bash
go run ./venue-service/cmd/seed/main.go
```

## Доступные Make-команды (пример для reservation-service)

```bash
make run      # Запуск reservation-service (из каталога reservation-service)
make dev      # Локальный запуск с hot-reload (при наличии конфигурации air)
make lint     # Проверка кода линтером
make fmt      # Форматирование кода
make vet      # Статический анализ кода
make test     # Запуск тестов (если добавлены)
```

## Команда проекта

[![tsuruevimran17](https://img.shields.io/badge/tsuruevimran17-181717?style=for-the-badge&logo=github&logoColor=white)](https://github.com/tsuruevimran17)
[![Idigov](https://img.shields.io/badge/Idigov_(Baisangur)-181717?style=for-the-badge&logo=github&logoColor=white)](https://github.com/Idigov)
[![namexamz](https://img.shields.io/badge/namexamz_(Хамзат)-181717?style=for-the-badge&logo=github&logoColor=white)](https://github.com/namexamz)
[![DjMariarty](https://img.shields.io/badge/DjMariarty-181717?style=for-the-badge&logo=github&logoColor=white)](https://github.com/DjMariarty)
