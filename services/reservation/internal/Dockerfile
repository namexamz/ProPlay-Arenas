FROM golang:1.25-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Сборка бинарника из директории с main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/rent-tools


FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app ./
EXPOSE 8081
CMD ["./app"]
