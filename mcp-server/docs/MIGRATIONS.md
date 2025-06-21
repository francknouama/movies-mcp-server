# Database Migrations Guide

## Overview
This project uses the [golang-migrate](https://github.com/golang-migrate/migrate) CLI tool for database schema migrations instead of embedding migration logic in the application code.

## Why CLI Tool Instead of Embedded Code?

### Advantages of CLI Approach:
1. **Separation of Concerns** - Database migrations are separate from application logic
2. **Standard Tool** - Uses the widely-adopted golang-migrate CLI
3. **Flexibility** - Can run migrations independently of the application
4. **CI/CD Friendly** - Easy to integrate into deployment pipelines
5. **No Runtime Dependencies** - Application doesn't need migration libraries
6. **Better DevEx** - Simpler debugging and management

### Previous Approach Issues:
- Required importing migration packages into the application
- Added runtime dependencies that weren't needed in production
- Mixed database management with application concerns
- Made the binary larger with unused migration code

## Installation

The CLI tool is automatically installed when needed:

```bash
# Automatic installation
make install-migrate

# Manual installation
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## Usage

### Common Commands

```bash
# Run all pending migrations
make db-migrate

# Rollback one migration
make db-migrate-down

# Reset database (drop all tables and re-run)
make db-migrate-reset

# Check current migration version
make db-migrate-version

# Seed database with sample data
make db-seed
```

### Manual CLI Usage

```bash
# Set database URL
export DATABASE_URL="postgres://movies_user:movies_password@localhost:5432/movies_mcp?sslmode=disable"

# Run migrations
migrate -path migrations -database "$DATABASE_URL" up

# Rollback one step
migrate -path migrations -database "$DATABASE_URL" down 1

# Check version
migrate -path migrations -database "$DATABASE_URL" version

# Force version (if needed)
migrate -path migrations -database "$DATABASE_URL" force 1
```

## Migration Files

Migration files are located in the `migrations/` directory:

```
migrations/
├── 001_create_movies_table.up.sql    # Create movies table
├── 001_create_movies_table.down.sql  # Drop movies table
├── 002_add_indexes.up.sql             # Add performance indexes
└── 002_add_indexes.down.sql           # Drop performance indexes
```

### Naming Convention
- `{version}_{description}.{direction}.sql`
- Version: Sequential number (001, 002, etc.)
- Description: Brief description with underscores
- Direction: `up` for forward migration, `down` for rollback

## Creating New Migrations

```bash
# Create a new migration
migrate create -ext sql -dir migrations -seq add_user_ratings

# This creates:
# migrations/003_add_user_ratings.up.sql
# migrations/003_add_user_ratings.down.sql
```

## Integration with Application

The application code is now cleaner:

```go
// Before: Required migration dependencies
import (
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

// After: Only database driver needed
import (
    _ "github.com/lib/pq"
)
```

## Development Workflow

1. **Start Database**: `make docker-compose-up`
2. **Setup Database**: `make db-setup`
3. **Run Migrations**: `make db-migrate`
4. **Seed Data**: `make db-seed`
5. **Start Development**: `make build && make run`

## Production Deployment

In production, run migrations before starting the application:

```bash
# 1. Run migrations
make db-migrate

# 2. Start application
./movies-server
```

## Troubleshooting

### Migration Stuck
```bash
# Check current state
make db-migrate-version

# Force to specific version if needed
migrate -path migrations -database "$DATABASE_URL" force 1
```

### Reset Development Database
```bash
make db-migrate-reset
make db-seed
```

### Check Migration History
```bash
# Connect to database
psql $DATABASE_URL

# Check schema_migrations table
SELECT * FROM schema_migrations;
```

## Best Practices

1. **Always create both up and down migrations**
2. **Test migrations on a copy of production data**
3. **Keep migrations small and focused**
4. **Never modify existing migration files**
5. **Backup database before running migrations in production**
6. **Run migrations in a transaction when possible**