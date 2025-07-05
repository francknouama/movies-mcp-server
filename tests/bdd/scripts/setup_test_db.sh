#!/bin/bash
set -e

# BDD Test Database Setup Script
# This script sets up the test database for BDD scenarios

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default test database configuration
DEFAULT_DB_HOST="localhost"
DEFAULT_DB_PORT="5432"
DEFAULT_DB_NAME="movies_mcp_test"
DEFAULT_DB_USER="movies_user"
DEFAULT_DB_PASSWORD="movies_password"

# Use environment variables or defaults
DB_HOST=${DB_HOST:-$DEFAULT_DB_HOST}
DB_PORT=${DB_PORT:-$DEFAULT_DB_PORT}
DB_NAME=${DB_NAME:-$DEFAULT_DB_NAME}
DB_USER=${DB_USER:-$DEFAULT_DB_USER}
DB_PASSWORD=${DB_PASSWORD:-$DEFAULT_DB_PASSWORD}

echo -e "${GREEN}Setting up BDD test database...${NC}"
echo -e "Host: ${YELLOW}$DB_HOST:$DB_PORT${NC}"
echo -e "Database: ${YELLOW}$DB_NAME${NC}"
echo -e "User: ${YELLOW}$DB_USER${NC}"

# Function to check if PostgreSQL is available
check_postgres() {
    if command -v psql &> /dev/null; then
        return 0
    elif command -v docker &> /dev/null && docker ps | grep -q postgres; then
        return 0
    else
        return 1
    fi
}

# Function to execute SQL with docker or local psql
exec_sql() {
    local sql="$1"
    local database="${2:-postgres}"
    
    if command -v docker &> /dev/null && docker ps | grep -q postgres; then
        # Use Docker PostgreSQL
        local container_id=$(docker ps | grep postgres | awk '{print $1}')
        echo -e "${YELLOW}Executing SQL in Docker container: $container_id${NC}"
        docker exec "$container_id" psql -U postgres -d "$database" -c "$sql"
    elif command -v psql &> /dev/null; then
        # Use local PostgreSQL
        echo -e "${YELLOW}Executing SQL on local PostgreSQL${NC}"
        PGPASSWORD=postgres psql -h "$DB_HOST" -p "$DB_PORT" -U postgres -d "$database" -c "$sql"
    else
        echo -e "${RED}No PostgreSQL client available${NC}"
        return 1
    fi
}

# Check if PostgreSQL is available
if ! check_postgres; then
    echo -e "${RED}Error: PostgreSQL is not available${NC}"
    echo "Please ensure either:"
    echo "1. PostgreSQL is installed locally with psql command"
    echo "2. PostgreSQL is running in Docker"
    exit 1
fi

echo -e "${GREEN}Creating test database and user...${NC}"

# Create test database
echo "Creating database $DB_NAME..."
exec_sql "CREATE DATABASE $DB_NAME;" postgres 2>/dev/null || {
    echo -e "${YELLOW}Database $DB_NAME may already exist${NC}"
}

# Create test user
echo "Creating user $DB_USER..."
exec_sql "CREATE USER $DB_USER WITH PASSWORD '$DB_PASSWORD';" postgres 2>/dev/null || {
    echo -e "${YELLOW}User $DB_USER may already exist${NC}"
}

# Grant privileges
echo "Granting privileges..."
exec_sql "GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;" postgres
exec_sql "ALTER DATABASE $DB_NAME OWNER TO $DB_USER;" postgres

# Connect to test database and grant schema privileges
echo "Setting up schema privileges..."
exec_sql "GRANT ALL ON SCHEMA public TO $DB_USER;" "$DB_NAME" 2>/dev/null || true
exec_sql "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO $DB_USER;" "$DB_NAME" 2>/dev/null || true
exec_sql "GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO $DB_USER;" "$DB_NAME" 2>/dev/null || true

# Check if migrations need to be run
echo -e "${GREEN}Checking database schema...${NC}"

# Run migrations if they haven't been run yet
MIGRATIONS_DIR="../../../migrations"
if [ -d "$MIGRATIONS_DIR" ]; then
    echo "Found migrations directory at $MIGRATIONS_DIR"
    
    # Check if migrate tool is available
    if command -v migrate &> /dev/null; then
        echo "Running database migrations..."
        DATABASE_URL="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"
        migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" up
        echo -e "${GREEN}Migrations completed${NC}"
    else
        echo -e "${YELLOW}Warning: migrate tool not found${NC}"
        echo "Please install golang-migrate or run migrations manually"
        echo "See: https://github.com/golang-migrate/migrate"
    fi
else
    echo -e "${YELLOW}Warning: migrations directory not found at $MIGRATIONS_DIR${NC}"
fi

# Verify test database setup
echo -e "${GREEN}Verifying test database setup...${NC}"

# Test connection with test user credentials
if command -v psql &> /dev/null; then
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT 'Test connection successful' AS status;" &>/dev/null
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Test database connection successful${NC}"
    else
        echo -e "${RED}✗ Test database connection failed${NC}"
        exit 1
    fi
else
    echo -e "${YELLOW}Skipping connection test (psql not available locally)${NC}"
fi

echo -e "${GREEN}✓ BDD test database setup completed successfully!${NC}"
echo ""
echo "You can now run BDD tests with the following environment:"
echo "export DB_HOST=$DB_HOST"
echo "export DB_PORT=$DB_PORT"
echo "export DB_NAME=$DB_NAME"
echo "export DB_USER=$DB_USER"
echo "export DB_PASSWORD=$DB_PASSWORD"