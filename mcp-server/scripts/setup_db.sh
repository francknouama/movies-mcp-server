#!/bin/bash
set -e

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Default values
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-movies_mcp}
DB_USER=${DB_USER:-movies_user}
DB_PASSWORD=${DB_PASSWORD:-movies_password}

echo "Setting up PostgreSQL database..."
echo "Host: $DB_HOST:$DB_PORT"
echo "Database: $DB_NAME"
echo "User: $DB_USER"

# Check if we're using Docker or local PostgreSQL
if command -v docker &> /dev/null && docker ps &> /dev/null; then
    # Check if PostgreSQL is running in Docker
    if docker ps | grep -q postgres; then
        echo "Using PostgreSQL in Docker..."
        
        # Create database and user using Docker
        docker exec $(docker ps | grep postgres | awk '{print $1}') psql -U postgres -c "CREATE DATABASE $DB_NAME;" 2>/dev/null || echo "Database may already exist"
        docker exec $(docker ps | grep postgres | awk '{print $1}') psql -U postgres -c "CREATE USER $DB_USER WITH PASSWORD '$DB_PASSWORD';" 2>/dev/null || echo "User may already exist"
        docker exec $(docker ps | grep postgres | awk '{print $1}') psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;"
        docker exec $(docker ps | grep postgres | awk '{print $1}') psql -U postgres -c "ALTER DATABASE $DB_NAME OWNER TO $DB_USER;"
    else
        echo "PostgreSQL container not found. Please run 'make docker-compose-up' first."
        exit 1
    fi
elif command -v psql &> /dev/null; then
    # Use local PostgreSQL
    echo "Using local PostgreSQL installation..."
    
    # Create database and user
    psql -h $DB_HOST -p $DB_PORT -U postgres -c "CREATE DATABASE $DB_NAME;" 2>/dev/null || echo "Database may already exist"
    psql -h $DB_HOST -p $DB_PORT -U postgres -c "CREATE USER $DB_USER WITH PASSWORD '$DB_PASSWORD';" 2>/dev/null || echo "User may already exist"
    psql -h $DB_HOST -p $DB_PORT -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;"
    psql -h $DB_HOST -p $DB_PORT -U postgres -c "ALTER DATABASE $DB_NAME OWNER TO $DB_USER;"
else
    echo "PostgreSQL client not found. Please install PostgreSQL or use Docker."
    exit 1
fi

echo "Database setup completed successfully!"