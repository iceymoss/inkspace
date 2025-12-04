.PHONY: help dev build run test docker-build docker-up docker-down clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

dev: ## Run in development mode
	go run main.go

build: ## Build the application
	go build -o bin/mysite main.go

run: build ## Build and run the application
	./bin/mysite

test: ## Run tests
	go test -v ./...

docker-build: ## Build Docker images
	docker-compose build

docker-up: ## Start Docker containers
	docker-compose up -d

docker-down: ## Stop Docker containers
	docker-compose down

docker-logs: ## Show Docker logs
	docker-compose logs -f

clean: ## Clean build files
	rm -rf bin/
	rm -rf uploads/*
	go clean

frontend-dev: ## Run frontend in development mode
	cd frontend && npm run dev

frontend-build: ## Build frontend
	cd frontend && npm run build

frontend-install: ## Install frontend dependencies
	cd frontend && npm install

