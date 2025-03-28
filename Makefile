build: ## Build the project
	@echo "Building the project..."
	go build -o main cmd/chat/main.go

start: ## Run the project
	go run cmd/chat/main.go

test: ## Run all tests
	go test ./...

db-run: ## Start the database using Docker
	docker-compose -f chatPostgres.docker-compose.yaml up

install: ## Install dependencies
	go mod tidy

swag-init:
		swag init -g internal/controllers/server.go -o api/openapi  --outputTypes "go,json,yaml" --overridesFile .swaggo

sqlc: ## Generate db schema
	sqlc generate

help: ## Display help message
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' Makefile
