.PHONY: help run fmt lint

help:
	@echo "Targets:"
	@echo "  run  - Run the service"
	@echo "  fmt  - Format Go code"
	@echo "  lint - Run linter"

run:
	go run cmd/app/main.go

fmt:
	go fmt ./...

lint:
	golangci-lint run
