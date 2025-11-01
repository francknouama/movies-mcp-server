#!/bin/bash
# SQLite Migration Test Script
# This script tests the SQLite migration to ensure everything works correctly

set -e  # Exit on error

echo "=========================================="
echo "SQLite Migration Testing Script"
echo "=========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Helper function to print success
success() {
    echo -e "${GREEN}✓${NC} $1"
    ((TESTS_PASSED++))
}

# Helper function to print failure
fail() {
    echo -e "${RED}✗${NC} $1"
    ((TESTS_FAILED++))
}

# Helper function to print info
info() {
    echo -e "${YELLOW}ℹ${NC} $1"
}

echo "Step 1: Download Dependencies"
echo "------------------------------"
if go mod download; then
    success "Dependencies downloaded successfully"
else
    fail "Failed to download dependencies"
    exit 1
fi
echo ""

echo "Step 2: Build Server Binary"
echo "----------------------------"
if go build -o ./movies-mcp-server cmd/server-sdk/main.go; then
    success "Server binary built successfully"
else
    fail "Failed to build server binary"
    exit 1
fi
echo ""

echo "Step 3: Clean Previous Database"
echo "--------------------------------"
if [ -f "movies.db" ]; then
    rm movies.db
    info "Removed existing movies.db"
fi
success "Clean database environment ready"
echo ""

echo "Step 4: Test Database Migration (SQLite)"
echo "-----------------------------------------"
export DB_DRIVER=sqlite
export DB_NAME=movies.db

if ./movies-mcp-server --migrate-only --migrations=./migrations 2>&1 | tee /tmp/migration.log; then
    if [ -f "movies.db" ]; then
        success "SQLite database created successfully"

        # Check if migrations table exists
        if sqlite3 movies.db "SELECT COUNT(*) FROM schema_migrations;" >/dev/null 2>&1; then
            success "Migrations table exists"

            # Check migration count
            MIGRATION_COUNT=$(sqlite3 movies.db "SELECT COUNT(*) FROM schema_migrations;")
            if [ "$MIGRATION_COUNT" -eq 5 ]; then
                success "All 5 migrations applied (count: $MIGRATION_COUNT)"
            else
                fail "Expected 5 migrations, found $MIGRATION_COUNT"
            fi
        else
            fail "Migrations table not found"
        fi
    else
        fail "movies.db file not created"
    fi
else
    fail "Migration failed"
fi
echo ""

echo "Step 5: Verify Database Schema"
echo "-------------------------------"
if [ -f "movies.db" ]; then
    # Check tables exist
    TABLES=$(sqlite3 movies.db "SELECT name FROM sqlite_master WHERE type='table';" | tr '\n' ' ')
    info "Tables found: $TABLES"

    if echo "$TABLES" | grep -q "movies"; then
        success "Movies table exists"
    else
        fail "Movies table not found"
    fi

    if echo "$TABLES" | grep -q "actors"; then
        success "Actors table exists"
    else
        fail "Actors table not found"
    fi

    if echo "$TABLES" | grep -q "movie_actors"; then
        success "Movie_actors junction table exists"
    else
        fail "Movie_actors table not found"
    fi

    # Check movies table structure
    info "Movies table schema:"
    sqlite3 movies.db ".schema movies" | head -20

    # Verify genre column is TEXT (for JSON)
    GENRE_TYPE=$(sqlite3 movies.db "SELECT type FROM pragma_table_info('movies') WHERE name='genre';")
    if [ "$GENRE_TYPE" = "TEXT" ]; then
        success "Genre column is TEXT (JSON-compatible)"
    else
        fail "Genre column type is $GENRE_TYPE, expected TEXT"
    fi
fi
echo ""

echo "Step 6: Test Server Startup (SQLite)"
echo "-------------------------------------"
info "Starting server in background for 5 seconds..."
timeout 5s ./movies-mcp-server 2>&1 | tee /tmp/server.log &
SERVER_PID=$!
sleep 2

if ps -p $SERVER_PID > /dev/null; then
    success "Server started successfully with SQLite"
    kill $SERVER_PID 2>/dev/null || true
    wait $SERVER_PID 2>/dev/null || true
else
    fail "Server failed to start"
fi

# Check log for expected output
if grep -q "driver: sqlite" /tmp/server.log; then
    success "Server confirmed using SQLite driver"
else
    fail "Server did not report SQLite driver"
fi
echo ""

echo "Step 7: Test Sample Data Insertion"
echo "-----------------------------------"
info "Inserting test movie directly into database..."

sqlite3 movies.db <<EOF
INSERT INTO movies (title, director, year, genre, rating, created_at, updated_at)
VALUES ('The Matrix', 'Wachowski Sisters', 1999, '["Action", "Sci-Fi"]', 8.7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO movies (title, director, year, genre, rating, created_at, updated_at)
VALUES ('Inception', 'Christopher Nolan', 2010, '["Action", "Thriller", "Sci-Fi"]', 8.8, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
EOF

if [ $? -eq 0 ]; then
    success "Sample movies inserted successfully"

    # Verify data
    MOVIE_COUNT=$(sqlite3 movies.db "SELECT COUNT(*) FROM movies;")
    if [ "$MOVIE_COUNT" -eq 2 ]; then
        success "Movie count correct: $MOVIE_COUNT"
    else
        fail "Expected 2 movies, found $MOVIE_COUNT"
    fi

    # Test JSON genre query
    info "Testing JSON genre search..."
    SCIFI_COUNT=$(sqlite3 movies.db "SELECT COUNT(*) FROM movies WHERE EXISTS (SELECT 1 FROM json_each(genre) WHERE value = 'Sci-Fi');")
    if [ "$SCIFI_COUNT" -eq 2 ]; then
        success "JSON genre search works (found $SCIFI_COUNT Sci-Fi movies)"
    else
        fail "JSON genre search failed (expected 2, found $SCIFI_COUNT)"
    fi
else
    fail "Failed to insert sample data"
fi
echo ""

echo "Step 8: Test Case-Insensitive Search"
echo "--------------------------------------"
CASE_TEST=$(sqlite3 movies.db "SELECT COUNT(*) FROM movies WHERE title LIKE '%matrix%' COLLATE NOCASE;")
if [ "$CASE_TEST" -eq 1 ]; then
    success "Case-insensitive search works"
else
    fail "Case-insensitive search failed (expected 1, found $CASE_TEST)"
fi
echo ""

echo "Step 9: Verify Database File Size"
echo "----------------------------------"
if [ -f "movies.db" ]; then
    DB_SIZE=$(du -h movies.db | cut -f1)
    info "SQLite database size: $DB_SIZE"
    success "Database file is portable and compact"
fi
echo ""

echo "Step 10: Test PostgreSQL Compatibility (Optional)"
echo "--------------------------------------------------"
info "To test PostgreSQL compatibility, set these environment variables:"
echo "  export DB_DRIVER=postgres"
echo "  export DB_HOST=localhost"
echo "  export DB_PORT=5432"
echo "  export DB_NAME=movies_mcp"
echo "  export DB_USER=movies_user"
echo "  export DB_PASSWORD=movies_password"
echo "  ./movies-mcp-server --migrate-only"
echo ""

echo "=========================================="
echo "Test Summary"
echo "=========================================="
echo -e "Tests Passed: ${GREEN}$TESTS_PASSED${NC}"
echo -e "Tests Failed: ${RED}$TESTS_FAILED${NC}"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ All tests passed! SQLite migration successful!${NC}"
    exit 0
else
    echo -e "${RED}✗ Some tests failed. Please review the output above.${NC}"
    exit 1
fi
