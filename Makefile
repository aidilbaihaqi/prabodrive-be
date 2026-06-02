# Auto-load .env so DB_* variables are available without exporting manually
ifneq (,$(wildcard .env))
    include .env
    export
endif

.PHONY: help run build test clean migrate-up migrate-down migrate-down-all migrate-fresh migrate-version migrate-create docker-build docker-up docker-down lint fmt

# ==========================================
# HELP
# ==========================================
help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

# ==========================================
# DEVELOPMENT
# ==========================================
run: ## Run the application
	go run ./cmd/api

dev: ## Run with hot reload (requires air)
	air

seed: ## Run database seeder (creates admin + sample user)
	go run ./cmd/seeder

build: ## Build the application
	go build -o bin/api ./cmd/api

clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf tmp/

# ==========================================
# TESTING
# ==========================================
test: ## Run all tests
	go test -v ./...

test-cover: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

test-unit: ## Run unit tests only
	go test -v -short ./...

test-integration: ## Run integration tests
	go test -v -tags=integration ./test/integration/...

# ==========================================
# DATABASE
# ==========================================
DB_URL ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

migrate-up: ## Run all pending migrations
	migrate -path migrations -database "$(DB_URL)" up

migrate-down: ## Rollback 1 migration
	migrate -path migrations -database "$(DB_URL)" down 1

migrate-down-all: ## Rollback semua migration
	migrate -path migrations -database "$(DB_URL)" down -all

migrate-fresh: ## Drop semua tabel lalu migrate ulang (hapus semua data!)
	migrate -path migrations -database "$(DB_URL)" down -all
	migrate -path migrations -database "$(DB_URL)" up

migrate-version: ## Cek versi migration saat ini
	migrate -path migrations -database "$(DB_URL)" version

migrate-create: ## Buat migration baru (usage: make migrate-create name=add_users_table)
	migrate create -ext sql -dir migrations -seq $(name)

# ==========================================
# DOCKER
# ==========================================
docker-build: ## Build Docker image
	docker build -t myapp:latest -f build/docker/Dockerfile .

docker-up: ## Start all containers
	docker-compose -f deployments/docker-compose.yml up -d

docker-down: ## Stop all containers
	docker-compose -f deployments/docker-compose.yml down

docker-logs: ## View container logs
	docker-compose -f deployments/docker-compose.yml logs -f

# ==========================================
# CODE QUALITY
# ==========================================
lint: ## Run linter
	golangci-lint run ./...

fmt: ## Format code
	go fmt ./...
	gofmt -s -w .

vet: ## Run go vet
	go vet ./...

# ==========================================
# DEPENDENCIES
# ==========================================
tidy: ## Tidy go modules
	go mod tidy

download: ## Download dependencies
	go mod download

vendor: ## Vendor dependencies
	go mod vendor

# ==========================================
# UTILITIES
# ==========================================
swagger: ## Generate Swagger documentation
	swag init -g cmd/api/main.go -o api/

mock: ## Generate mocks (requires mockgen)
	go generate ./...

# ==========================================
# PRODUCTION
# ==========================================
build-prod: ## Build for production
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o bin/api ./cmd/api
