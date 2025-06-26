#!/bin/bash

# IDM App Go - Setup Script
# This script helps set up the development environment

set -e

echo "🚀 Setting up IDM App Go development environment..."

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

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
print_status "Found Go version: $GO_VERSION"

# Check if Docker is installed (optional)
if command -v docker &> /dev/null; then
    print_status "Docker is available"
else
    print_warning "Docker is not installed. Docker Compose features will not work."
fi

# Check if PostgreSQL is installed (optional)
if command -v psql &> /dev/null; then
    print_status "PostgreSQL client is available"
else
    print_warning "PostgreSQL client is not installed. You'll need to use Docker for the database."
fi

# Install Go dependencies
print_status "Installing Go dependencies..."
go mod download
go mod tidy

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    print_status "Creating .env file from template..."
    cp env.example .env
    print_warning "Please edit .env file with your configuration"
else
    print_status ".env file already exists"
fi

# Create database if PostgreSQL is available
if command -v createdb &> /dev/null; then
    print_status "Creating database..."
    createdb idmapp 2>/dev/null || print_warning "Database 'idmapp' might already exist or you don't have permissions"
else
    print_warning "PostgreSQL client not found. You can create the database manually or use Docker Compose."
fi

# Build the application
print_status "Building the application..."
go build -o idmapp-go main.go

print_status "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "1. Edit .env file with your configuration"
echo "2. Run the application:"
echo "   - Local: ./idmapp-go or go run main.go"
echo "   - Docker: make docker-compose-up"
echo ""
echo "Available commands:"
echo "  make run              - Run the application"
echo "  make test             - Run tests"
echo "  make build            - Build the application"
echo "  make docker-compose-up   - Start with Docker Compose"
echo "  make help             - Show all available commands" 