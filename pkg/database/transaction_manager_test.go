package database

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	_ "modernc.org/sqlite"
)

func setupTransactionTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Create a test table
	schema := `
		CREATE TABLE tx_test_entities (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			value INTEGER
		);
	`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("failed to create test schema: %v", err)
	}

	return db
}

func TestNewTransactionManager(t *testing.T) {
	db := setupTransactionTestDB(t)
	defer db.Close()

	tm := NewTransactionManager(db)

	if tm == nil {
		t.Fatal("NewTransactionManager() returned nil")
	}

	if tm.db != db {
		t.Error("TransactionManager.db is not set correctly")
	}
}

func TestTransactionManager_WithTransaction_Success(t *testing.T) {
	db := setupTransactionTestDB(t)
	defer db.Close()

	tm := NewTransactionManager(db)
	ctx := context.Background()

	err := tm.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.Exec("INSERT INTO tx_test_entities (name, value) VALUES (?, ?)", "test", 42)
		return err
	})

	if err != nil {
		t.Fatalf("WithTransaction() error = %v", err)
	}

	// Verify the data was committed
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM tx_test_entities").Scan(&count)
	if err != nil {
		t.Fatalf("failed to verify transaction: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 row after transaction, got %d", count)
	}
}

func TestTransactionManager_WithTransaction_Rollback(t *testing.T) {
	db := setupTransactionTestDB(t)
	defer db.Close()

	tm := NewTransactionManager(db)
	ctx := context.Background()

	testErr := errors.New("test error")

	err := tm.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.Exec("INSERT INTO tx_test_entities (name, value) VALUES (?, ?)", "test", 42)
		if err != nil {
			return err
		}
		// Return error to trigger rollback
		return testErr
	})

	if err != testErr {
		t.Fatalf("WithTransaction() error = %v, want %v", err, testErr)
	}

	// Verify the data was rolled back
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM tx_test_entities").Scan(&count)
	if err != nil {
		t.Fatalf("failed to verify rollback: %v", err)
	}

	if count != 0 {
		t.Errorf("Expected 0 rows after rollback, got %d", count)
	}
}

func TestTransactionManager_WithTransaction_Panic(t *testing.T) {
	db := setupTransactionTestDB(t)
	defer db.Close()

	tm := NewTransactionManager(db)
	ctx := context.Background()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic to be re-thrown")
		}

		// Verify the data was rolled back
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM tx_test_entities").Scan(&count)
		if err != nil {
			t.Fatalf("failed to verify rollback after panic: %v", err)
		}

		if count != 0 {
			t.Errorf("Expected 0 rows after panic rollback, got %d", count)
		}
	}()

	_ = tm.WithTransaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.Exec("INSERT INTO tx_test_entities (name, value) VALUES (?, ?)", "test", 42)
		if err != nil {
			return err
		}
		// Trigger panic
		panic("test panic")
	})
}

func TestTransactionManager_WithTransaction_MultipleOperations(t *testing.T) {
	db := setupTransactionTestDB(t)
	defer db.Close()

	tm := NewTransactionManager(db)
	ctx := context.Background()

	err := tm.WithTransaction(ctx, func(tx *sql.Tx) error {
		// Insert multiple records
		for i := 1; i <= 3; i++ {
			_, err := tx.Exec("INSERT INTO tx_test_entities (name, value) VALUES (?, ?)", "test", i*10)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		t.Fatalf("WithTransaction() error = %v", err)
	}

	// Verify all data was committed
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM tx_test_entities").Scan(&count)
	if err != nil {
		t.Fatalf("failed to verify transaction: %v", err)
	}

	if count != 3 {
		t.Errorf("Expected 3 rows after transaction, got %d", count)
	}
}

func TestTransactionManager_WithTransactionResult_Success(t *testing.T) {
	db := setupTransactionTestDB(t)
	defer db.Close()

	tm := NewTransactionManager(db)
	ctx := context.Background()

	result, err := tm.WithTransactionResult(ctx, func(tx *sql.Tx) (interface{}, error) {
		_, err := tx.Exec("INSERT INTO tx_test_entities (name, value) VALUES (?, ?) RETURNING id", "test", 42)
		if err != nil {
			return nil, err
		}
		return "success", nil
	})

	if err != nil {
		t.Fatalf("WithTransactionResult() error = %v", err)
	}

	if result != "success" {
		t.Errorf("WithTransactionResult() result = %v, want success", result)
	}

	// Verify the data was committed
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM tx_test_entities").Scan(&count)
	if err != nil {
		t.Fatalf("failed to verify transaction: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 row after transaction, got %d", count)
	}
}

func TestTransactionManager_WithTransactionResult_Error(t *testing.T) {
	db := setupTransactionTestDB(t)
	defer db.Close()

	tm := NewTransactionManager(db)
	ctx := context.Background()

	testErr := errors.New("test error")

	result, err := tm.WithTransactionResult(ctx, func(tx *sql.Tx) (interface{}, error) {
		_, err := tx.Exec("INSERT INTO tx_test_entities (name, value) VALUES (?, ?)", "test", 42)
		if err != nil {
			return nil, err
		}
		return nil, testErr
	})

	if err != testErr {
		t.Fatalf("WithTransactionResult() error = %v, want %v", err, testErr)
	}

	if result != nil {
		t.Errorf("WithTransactionResult() result = %v, want nil", result)
	}

	// Verify the data was rolled back
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM tx_test_entities").Scan(&count)
	if err != nil {
		t.Fatalf("failed to verify rollback: %v", err)
	}

	if count != 0 {
		t.Errorf("Expected 0 rows after rollback, got %d", count)
	}
}

func TestNewTransactionHelper(t *testing.T) {
	db := setupTransactionTestDB(t)
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	helper := NewTransactionHelper(tx)

	if helper == nil {
		t.Fatal("NewTransactionHelper() returned nil")
	}

	if helper.tx != tx {
		t.Error("TransactionHelper.tx is not set correctly")
	}
}

func TestTransactionHelper_ExecContext(t *testing.T) {
	db := setupTransactionTestDB(t)
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	helper := NewTransactionHelper(tx)
	ctx := context.Background()

	result, err := helper.ExecContext(ctx, "INSERT INTO tx_test_entities (name, value) VALUES (?, ?)", "test", 42)
	if err != nil {
		t.Fatalf("ExecContext() error = %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("RowsAffected() error = %v", err)
	}

	if rowsAffected != 1 {
		t.Errorf("RowsAffected() = %d, want 1", rowsAffected)
	}
}

func TestTransactionHelper_QueryRowContext(t *testing.T) {
	db := setupTransactionTestDB(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec("INSERT INTO tx_test_entities (name, value) VALUES (?, ?)", "test", 42)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	helper := NewTransactionHelper(tx)
	ctx := context.Background()

	row := helper.QueryRowContext(ctx, "SELECT name, value FROM tx_test_entities WHERE id = ?", 1)

	var name string
	var value int
	if err := row.Scan(&name, &value); err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	if name != "test" || value != 42 {
		t.Errorf("QueryRowContext() got name=%s, value=%d, want name=test, value=42", name, value)
	}
}

func TestTransactionHelper_QueryContext(t *testing.T) {
	db := setupTransactionTestDB(t)
	defer db.Close()

	// Insert test data
	for i := 1; i <= 3; i++ {
		_, err := db.Exec("INSERT INTO tx_test_entities (name, value) VALUES (?, ?)", "test", i*10)
		if err != nil {
			t.Fatalf("failed to insert test data: %v", err)
		}
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	helper := NewTransactionHelper(tx)
	ctx := context.Background()

	rows, err := helper.QueryContext(ctx, "SELECT id, name, value FROM tx_test_entities ORDER BY id")
	if err != nil {
		t.Fatalf("QueryContext() error = %v", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id, value int
		var name string
		if err := rows.Scan(&id, &name, &value); err != nil {
			t.Fatalf("Scan() error = %v", err)
		}
		count++
	}

	if count != 3 {
		t.Errorf("QueryContext() returned %d rows, want 3", count)
	}
}

func TestTransactionHelper_CheckRowsAffected(t *testing.T) {
	db := setupTransactionTestDB(t)
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	helper := NewTransactionHelper(tx)
	ctx := context.Background()

	tests := []struct {
		name       string
		query      string
		args       []interface{}
		entityType string
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "rows affected - success",
			query:      "INSERT INTO tx_test_entities (name, value) VALUES (?, ?)",
			args:       []interface{}{"test", 42},
			entityType: "test_entity",
			wantErr:    false,
		},
		{
			name:       "no rows affected - not found",
			query:      "UPDATE tx_test_entities SET name = ? WHERE id = ?",
			args:       []interface{}{"updated", 999},
			entityType: "test_entity",
			wantErr:    true,
			errMsg:     "test_entity not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := helper.ExecContext(ctx, tt.query, tt.args...)
			if err != nil {
				t.Fatalf("ExecContext() error = %v", err)
			}

			err = helper.CheckRowsAffected(result, tt.entityType)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckRowsAffected() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("CheckRowsAffected() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestTransactionHelper_InsertWithID(t *testing.T) {
	db := setupTransactionTestDB(t)
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	helper := NewTransactionHelper(tx)
	ctx := context.Background()

	query := "INSERT INTO tx_test_entities (name, value) VALUES (?, ?) RETURNING id"
	id, err := helper.InsertWithID(ctx, query, "test", 42)
	if err != nil {
		t.Fatalf("InsertWithID() error = %v", err)
	}

	if id <= 0 {
		t.Errorf("InsertWithID() returned invalid id = %d", id)
	}

	// Verify within transaction
	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM tx_test_entities WHERE id = ?", id).Scan(&count)
	if err != nil {
		t.Fatalf("failed to verify insert: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 row with id=%d, got %d", id, count)
	}
}

func TestTransactionHelper_Update(t *testing.T) {
	db := setupTransactionTestDB(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec("INSERT INTO tx_test_entities (name, value) VALUES (?, ?)", "test", 42)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	helper := NewTransactionHelper(tx)
	ctx := context.Background()

	tests := []struct {
		name       string
		query      string
		entityType string
		args       []interface{}
		wantErr    bool
	}{
		{
			name:       "update existing entity",
			query:      "UPDATE tx_test_entities SET name = ? WHERE id = ?",
			entityType: "test_entity",
			args:       []interface{}{"updated", 1},
			wantErr:    false,
		},
		{
			name:       "update non-existent entity",
			query:      "UPDATE tx_test_entities SET name = ? WHERE id = ?",
			entityType: "test_entity",
			args:       []interface{}{"updated", 999},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := helper.Update(ctx, tt.query, tt.entityType, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransactionHelper_Delete(t *testing.T) {
	db := setupTransactionTestDB(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec("INSERT INTO tx_test_entities (name, value) VALUES (?, ?)", "test", 42)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	helper := NewTransactionHelper(tx)
	ctx := context.Background()

	tests := []struct {
		name       string
		query      string
		entityType string
		args       []interface{}
		wantErr    bool
	}{
		{
			name:       "delete existing entity",
			query:      "DELETE FROM tx_test_entities WHERE id = ?",
			entityType: "test_entity",
			args:       []interface{}{1},
			wantErr:    false,
		},
		{
			name:       "delete non-existent entity",
			query:      "DELETE FROM tx_test_entities WHERE id = ?",
			entityType: "test_entity",
			args:       []interface{}{999},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := helper.Delete(ctx, tt.query, tt.entityType, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
