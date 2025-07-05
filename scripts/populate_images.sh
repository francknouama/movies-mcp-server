#!/bin/bash

# Movies MCP Server - Populate Images Script
# This script downloads and stores movie poster images in the database

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

# Load environment variables
if [ -f "$PROJECT_DIR/.env" ]; then
    log_info "Loading environment variables from .env"
    export $(cat "$PROJECT_DIR/.env" | grep -v '^#' | xargs)
else
    log_warning "No .env file found, using environment defaults"
fi

# Check if database is accessible
log_info "Testing database connection..."
if ! psql -h "${DB_HOST:-127.0.0.1}" -p "${DB_PORT:-5432}" -U "${DB_USER:-movies_user}" -d "${DB_NAME:-movies_mcp}" -c "SELECT 1;" > /dev/null 2>&1; then
    log_error "Cannot connect to database"
    log_error "Please ensure PostgreSQL is running and database credentials are correct"
    exit 1
fi
log_success "Database connection successful"

# Change to project directory
cd "$PROJECT_DIR"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    log_error "Go is not installed or not in PATH"
    exit 1
fi

# Build and run the image population tool
log_info "Building image population tool..."
if ! go build -o build/populate-images ./cmd/populate-images; then
    log_error "Failed to build image population tool"
    exit 1
fi
log_success "Build successful"

log_info "Running image population tool..."
./build/populate-images

# Check exit code
if [ $? -eq 0 ]; then
    log_success "Image population completed successfully!"
    
    # Show summary
    TOTAL_WITH_IMAGES=$(psql -h "${DB_HOST:-127.0.0.1}" -p "${DB_PORT:-5432}" -U "${DB_USER:-movies_user}" -d "${DB_NAME:-movies_mcp}" -t -c "SELECT COUNT(*) FROM movies WHERE poster_data IS NOT NULL;" | xargs)
    TOTAL_MOVIES=$(psql -h "${DB_HOST:-127.0.0.1}" -p "${DB_PORT:-5432}" -U "${DB_USER:-movies_user}" -d "${DB_NAME:-movies_mcp}" -t -c "SELECT COUNT(*) FROM movies;" | xargs)
    
    log_info "Summary: $TOTAL_WITH_IMAGES out of $TOTAL_MOVIES movies now have poster images"
else
    log_error "Image population failed"
    exit 1
fi