package database

import (
	"context"
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Create a test table
	schema := `
		CREATE TABLE test_entities (
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

func TestNewBaseRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewBaseRepository(db)

	if repo == nil {
		t.Fatal("NewBaseRepository() returned nil")
	}

	if repo.DB() != db {
		t.Error("DB() did not return the correct database connection")
	}
}

func TestBaseRepository_DB(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewBaseRepository(db)
	gotDB := repo.DB()

	if gotDB != db {
		t.Errorf("DB() = %v, want %v", gotDB, db)
	}
}

func TestBaseRepository_ExecContext(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewBaseRepository(db)
	ctx := context.Background()

	query := "INSERT INTO test_entities (name, value) VALUES (?, ?)"
	result, err := repo.ExecContext(ctx, query, "test", 42)
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

func TestBaseRepository_QueryRowContext(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewBaseRepository(db)
	ctx := context.Background()

	// Insert test data
	_, err := db.Exec("INSERT INTO test_entities (name, value) VALUES (?, ?)", "test", 42)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Query the data
	query := "SELECT name, value FROM test_entities WHERE id = ?"
	row := repo.QueryRowContext(ctx, query, 1)

	var name string
	var value int
	if err := row.Scan(&name, &value); err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	if name != "test" || value != 42 {
		t.Errorf("QueryRowContext() got name=%s, value=%d, want name=test, value=42", name, value)
	}
}

func TestBaseRepository_QueryContext(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewBaseRepository(db)
	ctx := context.Background()

	// Insert test data
	for i := 1; i <= 3; i++ {
		_, err := db.Exec("INSERT INTO test_entities (name, value) VALUES (?, ?)", "test", i*10)
		if err != nil {
			t.Fatalf("failed to insert test data: %v", err)
		}
	}

	// Query multiple rows
	query := "SELECT id, name, value FROM test_entities ORDER BY id"
	rows, err := repo.QueryContext(ctx, query)
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

func TestBaseRepository_CheckRowsAffected(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewBaseRepository(db)
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
			query:      "INSERT INTO test_entities (name, value) VALUES (?, ?)",
			args:       []interface{}{"test", 42},
			entityType: "test_entity",
			wantErr:    false,
		},
		{
			name:       "no rows affected - not found",
			query:      "UPDATE test_entities SET name = ? WHERE id = ?",
			args:       []interface{}{"updated", 999},
			entityType: "test_entity",
			wantErr:    true,
			errMsg:     "test_entity not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.ExecContext(ctx, tt.query, tt.args...)
			if err != nil {
				t.Fatalf("ExecContext() error = %v", err)
			}

			err = repo.CheckRowsAffected(result, tt.entityType)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckRowsAffected() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("CheckRowsAffected() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestBaseRepository_Count(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewBaseRepository(db)
	ctx := context.Background()

	// Insert test data
	for i := 1; i <= 5; i++ {
		_, err := db.Exec("INSERT INTO test_entities (name, value) VALUES (?, ?)", "test", i)
		if err != nil {
			t.Fatalf("failed to insert test data: %v", err)
		}
	}

	tests := []struct {
		name      string
		query     string
		args      []interface{}
		wantCount int
		wantErr   bool
	}{
		{
			name:      "count all",
			query:     "SELECT COUNT(*) FROM test_entities",
			args:      nil,
			wantCount: 5,
			wantErr:   false,
		},
		{
			name:      "count with filter",
			query:     "SELECT COUNT(*) FROM test_entities WHERE value > ?",
			args:      []interface{}{3},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "count zero",
			query:     "SELECT COUNT(*) FROM test_entities WHERE value > ?",
			args:      []interface{}{100},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count, err := repo.Count(ctx, tt.query, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if count != tt.wantCount {
				t.Errorf("Count() = %d, want %d", count, tt.wantCount)
			}
		})
	}
}

func TestBaseRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewBaseRepository(db)
	ctx := context.Background()

	// Insert test data
	_, err := db.Exec("INSERT INTO test_entities (name, value) VALUES (?, ?)", "test", 42)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	tests := []struct {
		name       string
		query      string
		entityType string
		args       []interface{}
		wantErr    bool
	}{
		{
			name:       "delete existing entity",
			query:      "DELETE FROM test_entities WHERE id = ?",
			entityType: "test_entity",
			args:       []interface{}{1},
			wantErr:    false,
		},
		{
			name:       "delete non-existent entity",
			query:      "DELETE FROM test_entities WHERE id = ?",
			entityType: "test_entity",
			args:       []interface{}{999},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(ctx, tt.query, tt.entityType, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBaseRepository_InsertWithID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewBaseRepository(db)
	ctx := context.Background()

	query := "INSERT INTO test_entities (name, value) VALUES (?, ?) RETURNING id"
	id, err := repo.InsertWithID(ctx, query, "test", 42)
	if err != nil {
		t.Fatalf("InsertWithID() error = %v", err)
	}

	if id <= 0 {
		t.Errorf("InsertWithID() returned invalid id = %d", id)
	}

	// Verify the insert
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM test_entities WHERE id = ?", id).Scan(&count)
	if err != nil {
		t.Fatalf("failed to verify insert: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 row with id=%d, got %d", id, count)
	}
}

func TestBaseRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewBaseRepository(db)
	ctx := context.Background()

	// Insert test data
	_, err := db.Exec("INSERT INTO test_entities (name, value) VALUES (?, ?)", "test", 42)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	tests := []struct {
		name       string
		query      string
		entityType string
		args       []interface{}
		wantErr    bool
	}{
		{
			name:       "update existing entity",
			query:      "UPDATE test_entities SET name = ? WHERE id = ?",
			entityType: "test_entity",
			args:       []interface{}{"updated", 1},
			wantErr:    false,
		},
		{
			name:       "update non-existent entity",
			query:      "UPDATE test_entities SET name = ? WHERE id = ?",
			entityType: "test_entity",
			args:       []interface{}{"updated", 999},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(ctx, tt.query, tt.entityType, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBaseRepository_IsNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewBaseRepository(db)

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "sql.ErrNoRows",
			err:  sql.ErrNoRows,
			want: true,
		},
		{
			name: "other error",
			err:  sql.ErrConnDone,
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := repo.IsNotFound(tt.err)
			if got != tt.want {
				t.Errorf("IsNotFound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseRepository_WrapNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewBaseRepository(db)

	tests := []struct {
		name       string
		err        error
		entityType string
		wantErr    bool
		wantMsg    string
	}{
		{
			name:       "wrap sql.ErrNoRows",
			err:        sql.ErrNoRows,
			entityType: "test_entity",
			wantErr:    true,
			wantMsg:    "test_entity not found",
		},
		{
			name:       "other error - pass through",
			err:        sql.ErrConnDone,
			entityType: "test_entity",
			wantErr:    true,
			wantMsg:    "sql: connection is already closed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.WrapNotFound(tt.err, tt.entityType)
			if (err != nil) != tt.wantErr {
				t.Errorf("WrapNotFound() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && err.Error() != tt.wantMsg {
				t.Errorf("WrapNotFound() error message = %v, want %v", err.Error(), tt.wantMsg)
			}
		})
	}
}

func TestBaseRepository_Count_InvalidQuery(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewBaseRepository(db)
	ctx := context.Background()

	// Invalid query should return error
	_, err := repo.Count(ctx, "SELECT COUNT(*) FROM non_existent_table")
	if err == nil {
		t.Error("Count() should return error for invalid table")
	}
}

func TestBaseRepository_Delete_DatabaseError(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewBaseRepository(db)
	ctx := context.Background()

	// Invalid query should return error
	err := repo.Delete(ctx, "DELETE FROM non_existent_table WHERE id = ?", "test_entity", 1)
	if err == nil {
		t.Error("Delete() should return error for invalid table")
	}
}

func TestBaseRepository_Update_DatabaseError(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewBaseRepository(db)
	ctx := context.Background()

	// Invalid query should return error
	err := repo.Update(ctx, "UPDATE non_existent_table SET name = ? WHERE id = ?", "test_entity", "test", 1)
	if err == nil {
		t.Error("Update() should return error for invalid table")
	}
}

func TestBaseRepository_InsertWithID_InvalidQuery(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewBaseRepository(db)
	ctx := context.Background()

	// Invalid query should return error
	_, err := repo.InsertWithID(ctx, "INSERT INTO non_existent_table (name) VALUES (?) RETURNING id", "test")
	if err == nil {
		t.Error("InsertWithID() should return error for invalid table")
	}
}
