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
DB_SSLMODE=${DB_SSLMODE:-disable}

# Build connection string
DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

echo "Checking migration version..."
echo "Database: $DB_NAME on $DB_HOST:$DB_PORT"

# Check if migrate tool is installed, install if not
if ! command -v migrate &> /dev/null; then
    echo "Installing golang-migrate..."
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    
    # Check if it's in PATH
    if ! command -v migrate &> /dev/null; then
        echo "Adding GOPATH/bin to PATH for this session..."
        export PATH="$PATH:$(go env GOPATH)/bin"
    fi
fi

# Show current version
VERSION=$(migrate -path migrations -database "$DATABASE_URL" version 2>&1)
echo "Current migration version: $VERSION"

# Show available migrations
echo ""
echo "Available migrations:"
ls -la migrations/*.sql | awk '{print $9}' | sort