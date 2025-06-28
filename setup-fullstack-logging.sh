#!/bin/bash

# Fullstack IDM Setup with Centralized Logging
# This script sets up the complete IDM application with Go backend, React frontend, and ELK stack

set -e

echo "🚀 Starting fullstack IDM setup with centralized logging..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is running
check_docker() {
    print_status "Checking Docker..."
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker and try again."
        exit 1
    fi
    print_success "Docker is running"
}

# Clean up existing containers and volumes
cleanup() {
    print_status "Cleaning up existing containers and volumes..."
    
    # Stop and remove containers
    docker-compose -f docker-compose.fullstack-logging.yml down -v --remove-orphans 2>/dev/null || true
    
    # Remove any dangling images
    docker image prune -f 2>/dev/null || true
    
    # Remove any dangling volumes
    docker volume prune -f 2>/dev/null || true
    
    print_success "Cleanup completed"
}

# Check system resources
check_resources() {
    print_status "Checking system resources..."
    
    # Check available memory (in GB)
    if command -v sysctl >/dev/null 2>&1; then
        # macOS
        total_mem=$(sysctl -n hw.memsize | awk '{print $0/1024/1024/1024}')
    else
        # Linux
        total_mem=$(free -g | awk 'NR==2{print $2}')
    fi
    
    if (( $(echo "$total_mem < 8" | bc -l) )); then
        print_warning "System has less than 8GB RAM. Performance may be affected."
        print_warning "Consider closing other applications to free up memory."
    else
        print_success "System resources are adequate"
    fi
}

# Build and start services
start_services() {
    print_status "Building and starting services..."
    
    # Build images first
    print_status "Building Docker images..."
    docker-compose -f docker-compose.fullstack-logging.yml build --no-cache
    
    # Start services
    print_status "Starting services..."
    docker-compose -f docker-compose.fullstack-logging.yml up -d
    
    print_success "Services started successfully"
}

# Wait for services to be ready
wait_for_services() {
    print_status "Waiting for services to be ready..."
    
    # Wait for Elasticsearch
    print_status "Waiting for Elasticsearch..."
    timeout=120
    while ! curl -f http://localhost:9200/_cluster/health > /dev/null 2>&1; do
        if [ $timeout -le 0 ]; then
            print_error "Elasticsearch failed to start within 120 seconds"
            exit 1
        fi
        sleep 5
        timeout=$((timeout - 5))
    done
    print_success "Elasticsearch is ready"
    
    # Wait for Kibana
    print_status "Waiting for Kibana..."
    timeout=120
    while ! curl -f http://localhost:5601/api/status > /dev/null 2>&1; do
        if [ $timeout -le 0 ]; then
            print_error "Kibana failed to start within 120 seconds"
            exit 1
        fi
        sleep 5
        timeout=$((timeout - 5))
    done
    print_success "Kibana is ready"
    
    # Wait for Fluentd
    print_status "Waiting for Fluentd..."
    timeout=60
    while ! docker-compose -f docker-compose.fullstack-logging.yml exec -T fluentd pgrep -f fluentd > /dev/null 2>&1; do
        if [ $timeout -le 0 ]; then
            print_error "Fluentd failed to start within 60 seconds"
            exit 1
        fi
        sleep 2
        timeout=$((timeout - 2))
    done
    print_success "Fluentd is ready"
    
    # Wait for backend
    print_status "Waiting for Go backend..."
    timeout=60
    while ! curl -f http://localhost:8080/health > /dev/null 2>&1; do
        if [ $timeout -le 0 ]; then
            print_error "Go backend failed to start within 60 seconds"
            exit 1
        fi
        sleep 3
        timeout=$((timeout - 3))
    done
    print_success "Go backend is ready"
    
    # Wait for frontend
    print_status "Waiting for React frontend..."
    timeout=60
    while ! curl -f http://localhost:3000 > /dev/null 2>&1; do
        if [ $timeout -le 0 ]; then
            print_error "React frontend failed to start within 60 seconds"
            exit 1
        fi
        sleep 3
        timeout=$((timeout - 3))
    done
    print_success "React frontend is ready"
}

# Display service information
show_info() {
    echo ""
    echo "🎉 Fullstack IDM setup completed successfully!"
    echo ""
    echo "📋 Service Information:"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "🌐 React Frontend:     http://localhost:3000"
    echo "🔧 Go Backend API:     http://localhost:8080"
    echo "📊 Kibana Dashboard:   http://localhost:5601"
    echo "🔍 Elasticsearch:      http://localhost:9200"
    echo "📝 Fluentd:            localhost:24224"
    echo ""
    echo "📚 Quick Start:"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "1. Open http://localhost:3000 in your browser"
    echo "2. View logs in Kibana: http://localhost:5601"
    echo "3. Check service status: docker-compose -f docker-compose.fullstack-logging.yml ps"
    echo "4. View logs: docker-compose -f docker-compose.fullstack-logging.yml logs -f [service-name]"
    echo ""
    echo "📊 Logging Setup:"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "• Backend logs: idm-backend-logs-* index in Kibana"
    echo "• Frontend logs: idm-frontend-logs-* index in Kibana"
    echo "• Create index patterns in Kibana for both indices"
    echo ""
    echo "🛠️  Useful Commands:"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "• Stop services:     docker-compose -f docker-compose.fullstack-logging.yml down"
    echo "• Restart services:  docker-compose -f docker-compose.fullstack-logging.yml restart"
    echo "• View all logs:     docker-compose -f docker-compose.fullstack-logging.yml logs -f"
    echo "• Clean everything:  docker-compose -f docker-compose.fullstack-logging.yml down -v --remove-orphans"
    echo ""
}

# Main execution
main() {
    echo "🔧 Fullstack IDM Setup with Centralized Logging"
    echo "================================================"
    
    check_docker
    check_resources
    cleanup
    start_services
    wait_for_services
    show_info
}

# Run main function
main "$@" 