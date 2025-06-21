package database

import (
	"context"
	"database/sql"
	"fmt"
)

// TransactionManager provides utilities for transaction management
type TransactionManager struct {
	db *sql.DB
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(db *sql.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

// WithTransaction executes a function within a database transaction
// The transaction is automatically committed on success or rolled back on error
func (tm *TransactionManager) WithTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := tm.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	
	// Ensure rollback is called if we don't reach the commit
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-throw panic after rollback
		} else if err != nil {
			tx.Rollback()
		}
	}()
	
	// Execute the function with the transaction
	err = fn(tx)
	if err != nil {
		return err // Rollback will be called by defer
	}
	
	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return nil
}

// WithTransactionResult executes a function within a transaction and returns a result
// The transaction is automatically committed on success or rolled back on error
// Note: Using interface{} for broader compatibility - caller should type assert result
func (tm *TransactionManager) WithTransactionResult(ctx context.Context, fn func(*sql.Tx) (interface{}, error)) (interface{}, error) {
	var result interface{}
	var err error
	
	txErr := tm.WithTransaction(ctx, func(tx *sql.Tx) error {
		result, err = fn(tx)
		return err
	})
	
	if txErr != nil {
		return nil, txErr
	}
	
	return result, nil
}

// TransactionHelper provides common transaction operations
type TransactionHelper struct {
	tx *sql.Tx
}

// NewTransactionHelper creates a new transaction helper
func NewTransactionHelper(tx *sql.Tx) *TransactionHelper {
	return &TransactionHelper{tx: tx}
}

// ExecContext executes a query within the transaction
func (th *TransactionHelper) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return th.tx.ExecContext(ctx, query, args...)
}

// QueryRowContext executes a query that returns a single row within the transaction
func (th *TransactionHelper) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return th.tx.QueryRowContext(ctx, query, args...)
}

// QueryContext executes a query that returns multiple rows within the transaction
func (th *TransactionHelper) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return th.tx.QueryContext(ctx, query, args...)
}

// CheckRowsAffected validates that the expected number of rows were affected
func (th *TransactionHelper) CheckRowsAffected(result sql.Result, entityType string) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("%s not found", entityType)
	}
	
	return nil
}

// Insert executes an insert query and returns the new ID
func (th *TransactionHelper) InsertWithID(ctx context.Context, query string, args ...interface{}) (int, error) {
	var id int
	err := th.tx.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert record: %w", err)
	}
	return id, nil
}

// Update executes an update query and validates the result
func (th *TransactionHelper) Update(ctx context.Context, query string, entityType string, args ...interface{}) error {
	result, err := th.tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update %s: %w", entityType, err)
	}
	
	return th.CheckRowsAffected(result, entityType)
}

// Delete executes a delete query and validates the result
func (th *TransactionHelper) Delete(ctx context.Context, query string, entityType string, args ...interface{}) error {
	result, err := th.tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete %s: %w", entityType, err)
	}
	
	return th.CheckRowsAffected(result, entityType)
}