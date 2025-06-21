package database

import (
	"context"
	"database/sql"
	"fmt"
)

// BaseRepository provides common database operations for all repositories
type BaseRepository struct {
	db *sql.DB
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(db *sql.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

// DB returns the underlying database connection
func (r *BaseRepository) DB() *sql.DB {
	return r.db
}

// ExecContext executes a query with context and returns the result
func (r *BaseRepository) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return r.db.ExecContext(ctx, query, args...)
}

// QueryRowContext executes a query that returns a single row
func (r *BaseRepository) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return r.db.QueryRowContext(ctx, query, args...)
}

// QueryContext executes a query that returns multiple rows
func (r *BaseRepository) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return r.db.QueryContext(ctx, query, args...)
}

// CheckRowsAffected validates that the expected number of rows were affected
func (r *BaseRepository) CheckRowsAffected(result sql.Result, entityType string) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("%s not found", entityType)
	}
	
	return nil
}

// Count returns the count from a count query
func (r *BaseRepository) Count(ctx context.Context, query string, args ...interface{}) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to execute count query: %w", err)
	}
	return count, nil
}

// Delete executes a delete query and validates the result
func (r *BaseRepository) Delete(ctx context.Context, query string, entityType string, args ...interface{}) error {
	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete %s: %w", entityType, err)
	}
	
	return r.CheckRowsAffected(result, entityType)
}

// Insert executes an insert query and returns the new ID
func (r *BaseRepository) InsertWithID(ctx context.Context, query string, args ...interface{}) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert record: %w", err)
	}
	return id, nil
}

// Update executes an update query and validates the result
func (r *BaseRepository) Update(ctx context.Context, query string, entityType string, args ...interface{}) error {
	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update %s: %w", entityType, err)
	}
	
	return r.CheckRowsAffected(result, entityType)
}

// IsNotFound checks if an error is a "not found" error
func (r *BaseRepository) IsNotFound(err error) bool {
	return err == sql.ErrNoRows
}

// WrapNotFound wraps sql.ErrNoRows with a more descriptive error
func (r *BaseRepository) WrapNotFound(err error, entityType string) error {
	if err == sql.ErrNoRows {
		return fmt.Errorf("%s not found", entityType)
	}
	return err
}