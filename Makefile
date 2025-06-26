.PHONY: build run test clean docker-build docker-run docker-compose-up docker-compose-down

# Build the application
build:
	go build -o idmapp-go main.go

# Run the application
run:
	go run main.go

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -f idmapp-go

# Build Docker image
docker-build:
	docker build -t idmapp-go .

# Run Docker container
docker-run:
	docker run -p 8080:8080 --env-file .env idmapp-go

# Start services with docker-compose
docker-compose-up:
	docker-compose up -d

# Stop services with docker-compose
docker-compose-down:
	docker-compose down

# Show logs
logs:
	docker-compose logs -f app

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Install dependencies
deps:
	go mod download
	go mod tidy

# Generate go.sum
sum:
	go mod tidy
	go mod download

# Database migrations (auto-migrate with GORM)
migrate:
	go run main.go

# Help
help:
	@echo "Available commands:"
	@echo "  build           - Build the application"
	@echo "  run             - Run the application"
	@echo "  test            - Run tests"
	@echo "  clean           - Clean build artifacts"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-run      - Run Docker container"
	@echo "  docker-compose-up   - Start services with docker-compose"
	@echo "  docker-compose-down - Stop services with docker-compose"
	@echo "  logs            - Show application logs"
	@echo "  fmt             - Format code"
	@echo "  lint            - Run linter"
	@echo "  deps            - Install dependencies"
	@echo "  sum             - Generate go.sum"
	@echo "  migrate         - Run database migrations"
	@echo "  help            - Show this help" 