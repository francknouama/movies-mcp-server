#!/bin/bash

# Movies MCP Server - Database Seeding Script
# This script loads sample movie data into the database

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
Movies MCP Server Database Seeding Script

Usage: $0 [OPTIONS]

Options:
    -e, --env           Environment (development, staging, production) [default: $DEFAULT_ENV]
    -f, --force         Force seeding (clear existing data)
    -c, --clear-only    Only clear existing data (don't insert new data)
    -v, --verbose       Verbose output
    -h, --help          Show this help message

Database Connection:
    The script uses the following environment variables (in order of preference):
    1. DATABASE_URL - Full PostgreSQL connection string
    2. DB_* variables - Individual database connection parameters
    3. Default local development settings

Examples:
    $0                              # Seed development database
    $0 -e production -f             # Force seed production database
    $0 -c                           # Clear existing data only
    $0 -v                           # Verbose output

Environment Variables:
    DATABASE_URL                    # postgres://user:pass@host:port/db
    DB_HOST                         # localhost
    DB_PORT                         # 5432
    DB_NAME                         # movies_db
    DB_USER                         # movies_user
    DB_PASSWORD                     # password

EOF
}

get_db_connection() {
    local env="$1"
    
    # Check if DATABASE_URL is set
    if [[ -n "${DATABASE_URL:-}" ]]; then
        echo "$DATABASE_URL"
        return 0
    fi
    
    # Build connection string from individual components
    local host="${DB_HOST:-localhost}"
    local port="${DB_PORT:-5432}"
    local name="${DB_NAME:-movies_db}"
    local user="${DB_USER:-movies_user}"
    local password="${DB_PASSWORD:-password}"
    
    # Adjust database name based on environment
    case "$env" in
        development)
            name="${name}_dev"
            ;;
        staging)
            name="${name}_staging"
            ;;
        production)
            # Use name as is for production
            ;;
    esac
    
    echo "postgres://${user}:${password}@${host}:${port}/${name}"
}

check_database_connection() {
    local db_url="$1"
    
    log_info "Testing database connection..."
    
    if psql "$db_url" -c "SELECT 1;" > /dev/null 2>&1; then
        log_success "Database connection successful"
        return 0
    else
        log_error "Cannot connect to database: $db_url"
        log_info "Please check your database configuration and ensure the database is running."
        return 1
    fi
}

check_table_exists() {
    local db_url="$1"
    
    log_info "Checking if movies table exists..."
    
    local table_exists
    table_exists=$(psql "$db_url" -t -c "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'movies');" | xargs)
    
    if [[ "$table_exists" == "t" ]]; then
        log_success "Movies table exists"
        return 0
    else
        log_error "Movies table does not exist"
        log_info "Please run database migrations first:"
        log_info "  ./movies-mcp-server --migrate"
        log_info "  OR"
        log_info "  psql '$db_url' -f migrations/001_initial.sql"
        return 1
    fi
}

get_current_movie_count() {
    local db_url="$1"
    
    local count
    count=$(psql "$db_url" -t -c "SELECT COUNT(*) FROM movies;" | xargs)
    echo "$count"
}

clear_existing_data() {
    local db_url="$1"
    local force="$2"
    
    local current_count
    current_count=$(get_current_movie_count "$db_url")
    
    if [[ "$current_count" -gt 0 ]]; then
        log_warning "Database currently contains $current_count movies"
        
        if [[ "$force" != "true" ]]; then
            read -p "Do you want to clear existing data? (y/N): " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                log_info "Keeping existing data. New movies will be added."
                return 0
            fi
        fi
        
        log_info "Clearing existing movie data..."
        psql "$db_url" -c "TRUNCATE movies RESTART IDENTITY CASCADE;" > /dev/null
        log_success "Existing data cleared"
    else
        log_info "Database is empty, proceeding with seeding"
    fi
}

seed_database() {
    local db_url="$1"
    
    log_info "Inserting sample movie data..."
    
    # Execute the seed SQL file
    if psql "$db_url" -f "$SCRIPT_DIR/seed_data.sql" > /dev/null; then
        log_success "Sample data inserted successfully"
        
        # Display summary
        local final_count
        final_count=$(get_current_movie_count "$db_url")
        log_success "Database now contains $final_count movies"
        
        # Show some sample data
        log_info "Sample of inserted data:"
        psql "$db_url" -c "SELECT title, director, year, genre, rating FROM movies ORDER BY rating DESC LIMIT 5;" 2>/dev/null || true
        
    else
        log_error "Failed to insert sample data"
        return 1
    fi
}

validate_seeded_data() {
    local db_url="$1"
    
    log_info "Validating seeded data..."
    
    local validation_results
    validation_results=$(psql "$db_url" -t -c "
        SELECT 
            'Total movies: ' || COUNT(*) ||
            ', Years: ' || MIN(year) || '-' || MAX(year) ||
            ', Avg rating: ' || ROUND(AVG(rating), 2) ||
            ', Genres: ' || COUNT(DISTINCT genre) ||
            ', Directors: ' || COUNT(DISTINCT director)
        FROM movies;
    " | xargs)
    
    log_success "Validation complete: $validation_results"
    
    # Check for any data quality issues
    local issues
    issues=$(psql "$db_url" -t -c "
        SELECT COUNT(*) FROM movies 
        WHERE title IS NULL 
           OR title = '' 
           OR director IS NULL 
           OR director = '' 
           OR year IS NULL 
           OR rating IS NULL;
    " | xargs)
    
    if [[ "$issues" -gt 0 ]]; then
        log_warning "Found $issues movies with missing required data"
    else
        log_success "All movies have complete required data"
    fi
}

# Main execution
main() {
    local env="$DEFAULT_ENV"
    local force="false"
    local clear_only="false"
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
            -c|--clear-only)
                clear_only="true"
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
            *)
                log_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    if [[ "$verbose" == "true" ]]; then
        set -x
    fi
    
    # Load environment variables if .env file exists
    if [[ -f "$PROJECT_DIR/.env" ]]; then
        log_info "Loading environment variables from .env"
        set -a
        source "$PROJECT_DIR/.env"
        set +a
    fi
    
    # Get database connection
    local db_url
    db_url=$(get_db_connection "$env")
    
    log_info "Seeding $env database..."
    
    # Check dependencies
    if ! command -v psql &> /dev/null; then
        log_error "psql (PostgreSQL client) is not installed"
        log_info "Please install PostgreSQL client tools"
        exit 1
    fi
    
    # Check database connection
    if ! check_database_connection "$db_url"; then
        exit 1
    fi
    
    # Check if movies table exists
    if ! check_table_exists "$db_url"; then
        exit 1
    fi
    
    # Clear existing data if requested
    clear_existing_data "$db_url" "$force"
    
    # Exit if only clearing
    if [[ "$clear_only" == "true" ]]; then
        log_success "Data clearing completed"
        exit 0
    fi
    
    # Seed the database
    if seed_database "$db_url"; then
        validate_seeded_data "$db_url"
        log_success "Database seeding completed successfully!"
    else
        log_error "Database seeding failed"
        exit 1
    fi
}

# Trap errors
trap 'log_error "Seeding failed on line $LINENO"' ERR

# Check if script is being sourced or executed
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi