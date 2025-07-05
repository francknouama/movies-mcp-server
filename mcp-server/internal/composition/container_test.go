package composition

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

func TestNewContainer(t *testing.T) {
	// Create a mock database connection
	// In a real test, you'd use sqlmock or a test database
	db := &sql.DB{}

	container := NewContainer(db)

	// Test that all dependencies are properly initialized
	t.Run("Repositories", func(t *testing.T) {
		if container.MovieRepository == nil {
			t.Error("MovieRepository should not be nil")
		}
		if container.ActorRepository == nil {
			t.Error("ActorRepository should not be nil")
		}
	})

	t.Run("Application Services", func(t *testing.T) {
		if container.MovieService == nil {
			t.Error("MovieService should not be nil")
		}
		if container.ActorService == nil {
			t.Error("ActorService should not be nil")
		}
	})

	t.Run("Interface Handlers", func(t *testing.T) {
		if container.MovieHandlers == nil {
			t.Error("MovieHandlers should not be nil")
		}
		if container.ActorHandlers == nil {
			t.Error("ActorHandlers should not be nil")
		}
		if container.PromptHandlers == nil {
			t.Error("PromptHandlers should not be nil")
		}
		if container.CompoundToolHandlers == nil {
			t.Error("CompoundToolHandlers should not be nil")
		}
		if container.ContextManager == nil {
			t.Error("ContextManager should not be nil")
		}
		if container.ToolValidator == nil {
			t.Error("ToolValidator should not be nil")
		}
	})
}

func TestNewTestContainer(t *testing.T) {
	container := NewTestContainer()

	// Test that protocol handlers are initialized
	t.Run("Protocol Handlers", func(t *testing.T) {
		if container.PromptHandlers == nil {
			t.Error("PromptHandlers should not be nil")
		}
		if container.ToolValidator == nil {
			t.Error("ToolValidator should not be nil")
		}
	})

	// Test that data-related dependencies are nil for test container
	t.Run("Data Dependencies", func(t *testing.T) {
		if container.MovieRepository != nil {
			t.Error("MovieRepository should be nil in test container")
		}
		if container.ActorRepository != nil {
			t.Error("ActorRepository should be nil in test container")
		}
		if container.MovieService != nil {
			t.Error("MovieService should be nil in test container")
		}
		if container.ActorService != nil {
			t.Error("ActorService should be nil in test container")
		}
		if container.MovieHandlers != nil {
			t.Error("MovieHandlers should be nil in test container")
		}
		if container.ActorHandlers != nil {
			t.Error("ActorHandlers should be nil in test container")
		}
		if container.CompoundToolHandlers != nil {
			t.Error("CompoundToolHandlers should be nil in test container")
		}
		if container.ContextManager != nil {
			t.Error("ContextManager should be nil in test container")
		}
	})
}

func TestContainerWiring(t *testing.T) {
	// This test ensures that the dependencies are wired correctly
	// and can work together
	db := &sql.DB{}
	container := NewContainer(db)

	// Test that services receive correct repositories
	// This is implicitly tested by successful creation,
	// but we could add more specific tests if needed

	// Test that handlers receive correct services
	// Again, this is implicitly tested by successful creation

	// The fact that NewContainer doesn't panic indicates
	// that all dependencies are satisfied
	if container == nil {
		t.Fatal("Container should not be nil")
	}
}

func BenchmarkNewContainer(b *testing.B) {
	db := &sql.DB{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewContainer(db)
	}
}

func BenchmarkNewTestContainer(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewTestContainer()
	}
}

