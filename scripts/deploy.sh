#!/bin/bash

# Movies MCP Server Deployment Script
# This script handles deployment to various environments

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
DEFAULT_ENV="development"

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

show_help() {
    cat << EOF
Movies MCP Server Deployment Script

Usage: $0 [OPTIONS] COMMAND

Commands:
    docker      Deploy using Docker Compose
    build       Build Docker images
    test        Run deployment tests
    clean       Clean up resources
    status      Check deployment status

Options:
    -e, --env       Environment (development, staging, production) [default: $DEFAULT_ENV]
    -f, --force     Force deployment (skip confirmations)
    -v, --verbose   Verbose output
    -h, --help      Show this help message

Examples:
    $0 docker                           # Deploy using Docker Compose (development)
    $0 build -e production              # Build production Docker images
    $0 test -e staging                  # Test staging deployment
    $0 status                           # Check deployment status

Environment Configuration:
    development     Local Docker Compose deployment
    staging         Staging Docker deployment
    production      Production Docker deployment

EOF
}

check_dependencies() {
    local deps=("$@")
    local missing=()
    
    for dep in "${deps[@]}"; do
        if ! command -v "$dep" &> /dev/null; then
            missing+=("$dep")
        fi
    done
    
    if [[ ${#missing[@]} -gt 0 ]]; then
        log_error "Missing dependencies: ${missing[*]}"
        log_info "Please install the missing dependencies and try again."
        exit 1
    fi
}

validate_environment() {
    local env="$1"
    case "$env" in
        development|staging|production)
            return 0
            ;;
        *)
            log_error "Invalid environment: $env"
            log_info "Valid environments: development, staging, production"
            exit 1
            ;;
    esac
}

build_images() {
    local env="$1"
    
    log_info "Building Docker images for $env environment..."
    
    cd "$PROJECT_DIR"
    
    case "$env" in
        development)
            docker build -t movies-mcp-server:dev .
            ;;
        staging)
            docker build -t movies-mcp-server:staging .
            ;;
        production)
            docker build -f Dockerfile.production -t movies-mcp-server:latest .
            docker build -f Dockerfile.production -t movies-mcp-server:prod .
            ;;
    esac
    
    log_success "Docker images built successfully"
}

deploy_docker() {
    local env="$1"
    
    check_dependencies "docker" "docker-compose"
    
    log_info "Deploying with Docker Compose ($env environment)..."
    
    cd "$PROJECT_DIR"
    
    # Check if .env file exists
    if [[ ! -f ".env" ]]; then
        log_warning ".env file not found. Creating from template..."
        cp .env.example .env
        log_warning "Please edit .env file with your configuration before running again."
        exit 1
    fi
    
    # Build images if needed
    if [[ "$FORCE" == "true" ]] || ! docker images | grep -q "movies-mcp-server"; then
        build_images "$env"
    fi
    
    # Deploy with docker-compose
    case "$env" in
        development)
            docker-compose up -d
            ;;
        staging|production)
            docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
            ;;
    esac
    
    # Wait for services to be ready
    log_info "Waiting for services to be ready..."
    sleep 10
    
    # Check service health
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        log_success "Movies MCP Server is running and healthy"
        docker-compose ps
    else
        log_error "Movies MCP Server health check failed"
        docker-compose logs movies-mcp-server
        exit 1
    fi
}


run_tests() {
    local env="$1"
    
    log_info "Running deployment tests for $env environment..."
    
    case "$env" in
        development|staging)
            # Test local deployment
            if curl -f http://localhost:8080/health > /dev/null 2>&1; then
                log_success "Health check passed"
            else
                log_error "Health check failed"
                return 1
            fi
            
            if curl -f http://localhost:8080/ready > /dev/null 2>&1; then
                log_success "Readiness check passed"
            else
                log_error "Readiness check failed"
                return 1
            fi
            ;;
        production)
            # Test production deployment (adjust URL as needed)
            local prod_url="${PROD_URL:-https://movies-mcp.yourdomain.com}"
            if curl -f "$prod_url/health" > /dev/null 2>&1; then
                log_success "Production health check passed"
            else
                log_error "Production health check failed"
                return 1
            fi
            ;;
    esac
    
    log_success "All tests passed"
}

check_status() {
    local env="$1"
    
    log_info "Checking deployment status for $env environment..."
    
    if command -v docker-compose > /dev/null; then
        docker-compose ps
    else
        docker ps --filter "name=movies"
    fi
}

cleanup() {
    local env="$1"
    
    if [[ "$FORCE" != "true" ]]; then
        read -p "Are you sure you want to clean up $env deployment? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "Cleanup cancelled"
            exit 0
        fi
    fi
    
    log_info "Cleaning up $env deployment..."
    
    docker-compose down -v
    docker image prune -f
    
    log_success "Cleanup completed"
}

# Main execution
main() {
    local env="$DEFAULT_ENV"
    local command=""
    local force="false"
    local verbose="false"
    
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -e|--env)
                env="$2"
                shift 2
                ;;
            -f|--force)
                force="true"
                shift
                ;;
            -v|--verbose)
                verbose="true"
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            docker|build|test|clean|status)
                command="$1"
                shift
                ;;
            *)
                log_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # Set global variables
    FORCE="$force"
    VERBOSE="$verbose"
    
    if [[ "$verbose" == "true" ]]; then
        set -x
    fi
    
    # Validate inputs
    if [[ -z "$command" ]]; then
        log_error "No command specified"
        show_help
        exit 1
    fi
    
    validate_environment "$env"
    
    # Execute command
    case "$command" in
        docker)
            deploy_docker "$env"
            ;;
        build)
            build_images "$env"
            ;;
        test)
            run_tests "$env"
            ;;
        clean)
            cleanup "$env"
            ;;
        status)
            check_status "$env"
            ;;
        *)
            log_error "Unknown command: $command"
            show_help
            exit 1
            ;;
    esac
}

# Trap errors
trap 'log_error "Deployment failed on line $LINENO"' ERR

# Run main function
main "$@"