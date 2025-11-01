#!/bin/bash
# SQLite Database Inspector
# Quick tool to inspect the movies.db database

DB_FILE="${1:-movies.db}"

if [ ! -f "$DB_FILE" ]; then
    echo "Error: Database file '$DB_FILE' not found"
    echo "Usage: $0 [database_file]"
    exit 1
fi

echo "=========================================="
echo "SQLite Database Inspector"
echo "Database: $DB_FILE"
echo "=========================================="
echo ""

# Check if sqlite3 is installed
if ! command -v sqlite3 &> /dev/null; then
    echo "Error: sqlite3 command not found"
    echo "Please install SQLite3 CLI tools"
    exit 1
fi

# Database info
echo "📊 Database Information"
echo "----------------------"
SIZE=$(du -h "$DB_FILE" | cut -f1)
echo "File size: $SIZE"
echo ""

# List all tables
echo "📋 Tables"
echo "---------"
sqlite3 "$DB_FILE" "SELECT name FROM sqlite_master WHERE type='table' ORDER BY name;" | while read table; do
    COUNT=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM $table;")
    echo "  • $table ($COUNT rows)"
done
echo ""

# Migration status
echo "🔄 Migrations Applied"
echo "---------------------"
sqlite3 "$DB_FILE" "SELECT version, applied_at FROM schema_migrations ORDER BY version;" | while IFS='|' read version applied_at; do
    echo "  ✓ Migration $version - $applied_at"
done
echo ""

# Movies summary
MOVIE_COUNT=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM movies;" 2>/dev/null || echo "0")
if [ "$MOVIE_COUNT" -gt 0 ]; then
    echo "🎬 Movies Summary"
    echo "-----------------"
    echo "Total movies: $MOVIE_COUNT"

    # Top rated
    echo ""
    echo "Top 5 rated movies:"
    sqlite3 "$DB_FILE" -header -column "SELECT title, director, year, rating FROM movies WHERE rating IS NOT NULL ORDER BY rating DESC LIMIT 5;"

    # Genres
    echo ""
    echo "Sample genres (JSON):"
    sqlite3 "$DB_FILE" "SELECT title, genre FROM movies LIMIT 3;" | while IFS='|' read title genre; do
        echo "  • $title: $genre"
    done
fi
echo ""

# Actors summary
ACTOR_COUNT=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM actors;" 2>/dev/null || echo "0")
if [ "$ACTOR_COUNT" -gt 0 ]; then
    echo "🎭 Actors Summary"
    echo "-----------------"
    echo "Total actors: $ACTOR_COUNT"

    echo ""
    echo "Sample actors:"
    sqlite3 "$DB_FILE" -header -column "SELECT name, birth_year FROM actors LIMIT 5;"
fi
echo ""

# Relationships
RELATIONSHIP_COUNT=$(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM movie_actors;" 2>/dev/null || echo "0")
if [ "$RELATIONSHIP_COUNT" -gt 0 ]; then
    echo "🔗 Relationships"
    echo "----------------"
    echo "Total movie-actor links: $RELATIONSHIP_COUNT"
fi
echo ""

# Schema details
echo "🏗️  Schema Details"
echo "------------------"
echo "Movies table:"
sqlite3 "$DB_FILE" ".schema movies"
echo ""

# Interactive mode
echo "=========================================="
echo "💡 Tip: For interactive SQL, run:"
echo "   sqlite3 $DB_FILE"
echo ""
echo "Useful commands:"
echo "   .tables              - List all tables"
echo "   .schema movies       - Show movies table schema"
echo "   SELECT * FROM movies LIMIT 5;"
echo "   .exit                - Exit sqlite3"
echo "=========================================="
