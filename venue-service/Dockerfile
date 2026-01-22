# Multi-stage build для venue-service

# Этап сборки
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Устанавливаем git для зависимостей
RUN apk add --no-cache git

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Сборка бинарника из cmd/app/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/app

# Финальный этап
FROM alpine:latest

WORKDIR /app

# Копируем бинарник из этапа сборки
COPY --from=builder /app/app ./

# Открываем порт (по умолчанию 8080, настраивается через ENV)
EXPOSE 8082

# Запускаем приложение
CMD ["./app"]
