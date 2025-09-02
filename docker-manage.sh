#!/bin/bash

# Theater Booking System Docker Management Script
# Usage: ./docker-manage.sh [command]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
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

# Function to show usage
show_usage() {
    cat << EOF
🎭 Theater Booking System Docker Management Script

Usage: $0 [COMMAND]

Commands:
  start|up       Start all services in detached mode
  stop|down      Stop all services
  restart        Restart all services
  rebuild        Rebuild and start all services
  logs           Show logs for all services
  logs-backend   Show logs for backend service
  logs-frontend  Show logs for frontend service
  logs-redis     Show logs for Redis service
  logs-mysql     Show logs for MySQL service
  status         Show status of all services
  cleanup        Remove all containers, networks, and volumes
  reset          Complete reset - rebuild everything from scratch
  shell-backend  Open shell in backend container
  shell-frontend Open shell in frontend container
  shell-mysql    Open MySQL shell
  shell-redis    Open Redis CLI
  health         Check health status of all services
  help           Show this help message

Examples:
  $0 start           # Start all services
  $0 logs-backend    # Show backend logs
  $0 shell-mysql     # Open MySQL shell
  $0 reset           # Complete reset and rebuild

EOF
}

# Function to check if Docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker first."
        exit 1
    fi
}

# Function to check if docker-compose.yaml exists
check_compose_file() {
    if [ ! -f "docker-compose.yaml" ]; then
        print_error "docker-compose.yaml not found in current directory"
        exit 1
    fi
}

# Start services
start_services() {
    print_info "Starting Theater Booking System..."
    docker-compose up -d
    print_success "All services started successfully!"
    print_info "Frontend: http://localhost:3000"
    print_info "Backend API: http://localhost:8080"
    print_info "MySQL: localhost:3306"
    print_info "Redis: localhost:6379"
}

# Stop services
stop_services() {
    print_info "Stopping Theater Booking System..."
    docker-compose down
    print_success "All services stopped successfully!"
}

# Restart services
restart_services() {
    print_info "Restarting Theater Booking System..."
    docker-compose restart
    print_success "All services restarted successfully!"
}

# Rebuild and start services
rebuild_services() {
    print_info "Rebuilding Theater Booking System..."
    docker-compose down
    docker-compose build --no-cache
    docker-compose up -d
    print_success "All services rebuilt and started successfully!"
}

# Show logs
show_logs() {
    case "${2:-all}" in
        "backend")
            docker-compose logs -f theater-backend
            ;;
        "frontend")
            docker-compose logs -f theater-website
            ;;
        "redis")
            docker-compose logs -f redis-theater
            ;;
        "mysql")
            docker-compose logs -f mysql-theater
            ;;
        *)
            docker-compose logs -f
            ;;
    esac
}

# Show status
show_status() {
    print_info "Theater Booking System Status:"
    docker-compose ps
}

# Cleanup everything
cleanup_system() {
    print_warning "This will remove ALL containers, networks, and volumes!"
    read -p "Are you sure? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_info "Cleaning up Theater Booking System..."
        docker-compose down -v --remove-orphans
        docker system prune -f
        print_success "Cleanup completed!"
    else
        print_info "Cleanup cancelled."
    fi
}

# Complete reset
reset_system() {
    print_warning "This will perform a complete reset - rebuild everything from scratch!"
    read -p "Are you sure? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_info "Performing complete reset..."
        docker-compose down -v --remove-orphans
        docker system prune -f
        docker-compose build --no-cache
        docker-compose up -d
        print_success "Reset completed successfully!"
    else
        print_info "Reset cancelled."
    fi
}

# Open shells
open_shell() {
    case "$2" in
        "backend")
            docker-compose exec theater-backend sh
            ;;
        "frontend")
            docker-compose exec theater-website sh
            ;;
        "mysql")
            docker-compose exec mysql-theater mysql -u theater_user -ptheater_password theater_booking
            ;;
        "redis")
            docker-compose exec redis-theater redis-cli -a theater_redis_pass
            ;;
        *)
            print_error "Invalid shell target. Use: backend, frontend, mysql, redis"
            exit 1
            ;;
    esac
}

# Check health
check_health() {
    print_info "Checking health status..."
    echo ""
    
    # Check backend health
    print_info "Backend Health:"
    if curl -s http://localhost:8080/api/v1/health > /dev/null; then
        print_success "✓ Backend is healthy"
    else
        print_error "✗ Backend is not responding"
    fi
    
    # Check frontend
    print_info "Frontend Health:"
    if curl -s http://localhost:3000 > /dev/null; then
        print_success "✓ Frontend is healthy"
    else
        print_error "✗ Frontend is not responding"
    fi
    
    # Check containers
    print_info "Container Status:"
    docker-compose ps
}

# Main script logic
main() {
    # Check prerequisites
    check_docker
    check_compose_file
    
    # Parse command
    case "${1:-help}" in
        "start"|"up")
            start_services
            ;;
        "stop"|"down")
            stop_services
            ;;
        "restart")
            restart_services
            ;;
        "rebuild")
            rebuild_services
            ;;
        "logs")
            show_logs "$@"
            ;;
        "logs-backend")
            show_logs "logs" "backend"
            ;;
        "logs-frontend")
            show_logs "logs" "frontend"
            ;;
        "logs-redis")
            show_logs "logs" "redis"
            ;;
        "logs-mysql")
            show_logs "logs" "mysql"
            ;;
        "status")
            show_status
            ;;
        "cleanup")
            cleanup_system
            ;;
        "reset")
            reset_system
            ;;
        "shell-backend")
            open_shell "shell" "backend"
            ;;
        "shell-frontend")
            open_shell "shell" "frontend"
            ;;
        "shell-mysql")
            open_shell "shell" "mysql"
            ;;
        "shell-redis")
            open_shell "shell" "redis"
            ;;
        "health")
            check_health
            ;;
        "help"|"-h"|"--help")
            show_usage
            ;;
        *)
            print_error "Unknown command: $1"
            echo ""
            show_usage
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"
