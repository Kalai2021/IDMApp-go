#!/bin/bash

# IDM App Go - Deployment Script
# This script helps deploy the application to production

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if .env file exists
if [ ! -f .env ]; then
    print_error ".env file not found. Please create it from env.example"
    exit 1
fi

# Build the application
print_status "Building the application..."
go build -ldflags="-s -w" -o idmapp-go main.go

# Run tests
print_status "Running tests..."
go test ./...

# Check if Docker is available for containerized deployment
if command -v docker &> /dev/null; then
    print_status "Building Docker image..."
    docker build -t idmapp-go:latest .
    
    print_status "Docker image built successfully!"
    echo ""
    echo "To run the container:"
    echo "  docker run -p 8080:8080 --env-file .env idmapp-go:latest"
    echo ""
    echo "To run with docker-compose:"
    echo "  docker-compose up -d"
else
    print_warning "Docker not available. Using binary deployment."
fi

print_status "✅ Deployment preparation complete!"
echo ""
echo "Next steps:"
echo "1. Ensure your .env file is configured for production"
echo "2. Ensure your database is accessible"
echo "3. Run the application:"
echo "   - Binary: ./idmapp-go"
echo "   - Docker: docker run -p 8080:8080 --env-file .env idmapp-go:latest"
echo "   - Docker Compose: docker-compose up -d" 