#!/bin/bash
set -e

# BDD Test Runner Script
# This script runs all BDD tests with proper setup

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default configuration
DEFAULT_DB_HOST="localhost"
DEFAULT_DB_PORT="5432"
DEFAULT_DB_NAME="movies_mcp_test"
DEFAULT_DB_USER="movies_user"
DEFAULT_DB_PASSWORD="movies_password"

# Use environment variables or defaults
export DB_HOST=${DB_HOST:-$DEFAULT_DB_HOST}
export DB_PORT=${DB_PORT:-$DEFAULT_DB_PORT}
export DB_NAME=${DB_NAME:-$DEFAULT_DB_NAME}
export DB_USER=${DB_USER:-$DEFAULT_DB_USER}
export DB_PASSWORD=${DB_PASSWORD:-$DEFAULT_DB_PASSWORD}

echo -e "${BLUE}===========================================${NC}"
echo -e "${BLUE}         Movies MCP Server BDD Tests${NC}"
echo -e "${BLUE}===========================================${NC}"
echo ""

# Check if we're in the right directory
if [ ! -f "godog_test.go" ] && [ ! -d "features" ]; then
    echo -e "${RED}Error: Not in BDD tests directory${NC}"
    echo "Please run this script from mcp-server/tests/bdd/"
    exit 1
fi

# Function to check prerequisites
check_prerequisites() {
    echo -e "${GREEN}Checking prerequisites...${NC}"
    
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        echo -e "${RED}Error: Go is not installed${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ Go is installed${NC}"
    
    # Check if database is accessible
    if command -v psql &> /dev/null; then
        PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1;" &>/dev/null
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}✓ Test database is accessible${NC}"
        else
            echo -e "${YELLOW}Warning: Test database connection failed${NC}"
            echo "You may need to run: ./scripts/setup_test_db.sh"
        fi
    else
        echo -e "${YELLOW}Warning: psql not available for database check${NC}"
    fi
    
    # Check if main server binary exists
    if [ -f "../../main" ]; then
        echo -e "${GREEN}✓ Main server binary found${NC}"
    else
        echo -e "${YELLOW}Warning: Main server binary not found at ../../main${NC}"
        echo "You may need to build the server first: go build -o main"
    fi
}

# Function to run specific test tags
run_tests() {
    local tags="$1"
    local description="$2"
    
    echo ""
    echo -e "${BLUE}Running $description...${NC}"
    
    if [ -n "$tags" ]; then
        echo -e "${YELLOW}Tags: $tags${NC}"
        go test -tags="$tags" -v
    else
        echo -e "${YELLOW}Running all scenarios${NC}"
        go test -v
    fi
}

# Main execution
main() {
    # Check prerequisites
    check_prerequisites
    
    echo ""
    echo -e "${GREEN}Environment Configuration:${NC}"
    echo -e "Database: ${YELLOW}$DB_HOST:$DB_PORT/$DB_NAME${NC}"
    echo -e "User: ${YELLOW}$DB_USER${NC}"
    echo ""
    
    # Parse command line arguments
    case "${1:-all}" in
        "smoke")
            run_tests "@smoke" "Smoke Tests"
            ;;
        "movies")
            run_tests "@movies" "Movie Operations Tests"
            ;;
        "actors")
            run_tests "@actors" "Actor Operations Tests"
            ;;
        "mcp")
            run_tests "@mcp" "MCP Protocol Tests"
            ;;
        "search")
            run_tests "@search" "Search Tests"
            ;;
        "integration")
            run_tests "@integration" "Integration Tests"
            ;;
        "error")
            run_tests "@error-handling" "Error Handling Tests"
            ;;
        "all"|"")
            echo -e "${GREEN}Running all BDD scenarios...${NC}"
            run_tests "" "All BDD Tests"
            ;;
        "help"|"-h"|"--help")
            echo -e "${GREEN}Usage: $0 [test-type]${NC}"
            echo ""
            echo -e "${YELLOW}Available test types:${NC}"
            echo "  smoke       - Quick smoke tests"
            echo "  movies      - Movie operations tests"
            echo "  actors      - Actor operations tests"
            echo "  mcp         - MCP protocol tests"
            echo "  search      - Search functionality tests"
            echo "  integration - Integration workflow tests"
            echo "  error       - Error handling tests"
            echo "  all         - All tests (default)"
            echo ""
            echo -e "${YELLOW}Environment variables:${NC}"
            echo "  DB_HOST, DB_PORT, DB_NAME, DB_USER, DB_PASSWORD"
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown test type: $1${NC}"
            echo "Use '$0 help' for available options"
            exit 1
            ;;
    esac
    
    echo ""
    echo -e "${GREEN}BDD test execution completed!${NC}"
}

# Run main function with all arguments
main "$@"