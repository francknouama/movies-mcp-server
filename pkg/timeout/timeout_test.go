package timeout

import (
	"context"
	"errors"
	"testing"
	"time"

	"movies-mcp-server/pkg/logging"
)

func TestDefaultTimeoutConfig(t *testing.T) {
	config := DefaultTimeoutConfig()

	if config.RequestTimeout != 30*time.Second {
		t.Errorf("Expected RequestTimeout to be 30s, got %v", config.RequestTimeout)
	}
	if config.DatabaseTimeout != 10*time.Second {
		t.Errorf("Expected DatabaseTimeout to be 10s, got %v", config.DatabaseTimeout)
	}
	if config.ImageTimeout != 15*time.Second {
		t.Errorf("Expected ImageTimeout to be 15s, got %v", config.ImageTimeout)
	}
	if config.ShutdownTimeout != 30*time.Second {
		t.Errorf("Expected ShutdownTimeout to be 30s, got %v", config.ShutdownTimeout)
	}
	if config.HealthCheckTimeout != 5*time.Second {
		t.Errorf("Expected HealthCheckTimeout to be 5s, got %v", config.HealthCheckTimeout)
	}
}

func TestNewManager(t *testing.T) {
	logger := logging.New(logging.LevelInfo)
	config := DefaultTimeoutConfig()

	manager := NewManager(config, logger)
	if manager.config != config {
		t.Error("Manager should use provided config")
	}
	if manager.logger != logger {
		t.Error("Manager should use provided logger")
	}

	// Test with nil config
	manager2 := NewManager(nil, logger)
	if manager2.config == nil {
		t.Error("Manager should use default config when nil provided")
	}
}

func TestWithRequestTimeout(t *testing.T) {
	manager := NewManager(DefaultTimeoutConfig(), nil)
	ctx := context.Background()

	timeoutCtx, cancel := manager.WithRequestTimeout(ctx)
	defer cancel()

	deadline, ok := timeoutCtx.Deadline()
	if !ok {
		t.Error("Context should have deadline")
	}

	expectedDeadline := time.Now().Add(30 * time.Second)
	if deadline.After(expectedDeadline.Add(1*time.Second)) || deadline.Before(expectedDeadline.Add(-1*time.Second)) {
		t.Errorf("Deadline should be approximately %v, got %v", expectedDeadline, deadline)
	}
}

func TestWithDatabaseTimeout(t *testing.T) {
	manager := NewManager(DefaultTimeoutConfig(), nil)
	ctx := context.Background()

	timeoutCtx, cancel := manager.WithDatabaseTimeout(ctx)
	defer cancel()

	deadline, ok := timeoutCtx.Deadline()
	if !ok {
		t.Error("Context should have deadline")
	}

	expectedDeadline := time.Now().Add(10 * time.Second)
	if deadline.After(expectedDeadline.Add(1*time.Second)) || deadline.Before(expectedDeadline.Add(-1*time.Second)) {
		t.Errorf("Deadline should be approximately %v, got %v", expectedDeadline, deadline)
	}
}

func TestWithCustomTimeout(t *testing.T) {
	manager := NewManager(DefaultTimeoutConfig(), nil)
	ctx := context.Background()
	customTimeout := 5 * time.Second

	timeoutCtx, cancel := manager.WithCustomTimeout(ctx, customTimeout, "test")
	defer cancel()

	deadline, ok := timeoutCtx.Deadline()
	if !ok {
		t.Error("Context should have deadline")
	}

	expectedDeadline := time.Now().Add(customTimeout)
	if deadline.After(expectedDeadline.Add(1*time.Second)) || deadline.Before(expectedDeadline.Add(-1*time.Second)) {
		t.Errorf("Deadline should be approximately %v, got %v", expectedDeadline, deadline)
	}
}

func TestHandleTimeout(t *testing.T) {
	manager := NewManager(DefaultTimeoutConfig(), nil)

	tests := []struct {
		name          string
		ctx           context.Context
		operation     string
		err           error
		expectTimeout bool
		expectCancel  bool
	}{
		{
			name:      "no error",
			ctx:       context.Background(),
			operation: "test",
			err:       nil,
		},
		{
			name:          "timeout error",
			ctx:           createTimedOutContext(),
			operation:     "test",
			err:           context.DeadlineExceeded,
			expectTimeout: true,
		},
		{
			name:         "cancelled error",
			ctx:          createCancelledContext(),
			operation:    "test",
			err:          context.Canceled,
			expectCancel: true,
		},
		{
			name:      "other error",
			ctx:       context.Background(),
			operation: "test",
			err:       errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.HandleTimeout(tt.ctx, tt.operation, tt.err)

			if tt.expectTimeout {
				if err == nil {
					t.Error("Expected timeout error")
				}
				// Check if it's a timeout error by examining the error message
				if err.Error() == "" {
					t.Error("Expected non-empty error message for timeout")
				}
			} else if tt.expectCancel {
				if err == nil {
					t.Error("Expected cancellation error")
				}
			} else if tt.err == nil {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			} else {
				if err != tt.err {
					t.Errorf("Expected original error %v, got %v", tt.err, err)
				}
			}
		})
	}
}

func createTimedOutContext() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	time.Sleep(time.Millisecond) // Ensure context times out
	return ctx
}

func createCancelledContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

func TestRequestTimeoutMiddleware(t *testing.T) {
	config := &TimeoutConfig{
		RequestTimeout: 100 * time.Millisecond,
	}
	manager := NewManager(config, nil)
	middleware := NewRequestTimeoutMiddleware(manager)

	t.Run("handler completes in time", func(t *testing.T) {
		handler := func(ctx context.Context) error {
			time.Sleep(50 * time.Millisecond)
			return nil
		}

		wrappedHandler := middleware.WrapHandler(handler)
		err := wrappedHandler(context.Background())
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("handler times out", func(t *testing.T) {
		handler := func(ctx context.Context) error {
			select {
			case <-time.After(200 * time.Millisecond):
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		wrappedHandler := middleware.WrapHandler(handler)
		err := wrappedHandler(context.Background())
		if err == nil {
			t.Error("Expected timeout error")
		}
	})

	t.Run("handler returns error", func(t *testing.T) {
		expectedErr := errors.New("handler error")
		handler := func(ctx context.Context) error {
			return expectedErr
		}

		wrappedHandler := middleware.WrapHandler(handler)
		err := wrappedHandler(context.Background())
		if err == nil {
			t.Error("Expected error")
		}
	})

	t.Run("handler panics", func(t *testing.T) {
		handler := func(ctx context.Context) error {
			panic("test panic")
		}

		wrappedHandler := middleware.WrapHandler(handler)
		err := wrappedHandler(context.Background())
		if err == nil {
			t.Error("Expected panic to be converted to error")
		}
	})
}

func TestDatabaseTimeoutWrapper(t *testing.T) {
	config := &TimeoutConfig{
		DatabaseTimeout: 100 * time.Millisecond,
	}
	manager := NewManager(config, nil)
	wrapper := NewDatabaseTimeoutWrapper(manager)

	t.Run("operation completes in time", func(t *testing.T) {
		operation := func(ctx context.Context) error {
			time.Sleep(50 * time.Millisecond)
			return nil
		}

		wrappedOp := wrapper.WrapDatabaseOperation(operation)
		err := wrappedOp(context.Background())
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("operation times out", func(t *testing.T) {
		operation := func(ctx context.Context) error {
			select {
			case <-time.After(200 * time.Millisecond):
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		wrappedOp := wrapper.WrapDatabaseOperation(operation)
		err := wrappedOp(context.Background())
		if err == nil {
			t.Error("Expected timeout error")
		}
	})
}

func TestImageTimeoutWrapper(t *testing.T) {
	config := &TimeoutConfig{
		ImageTimeout: 100 * time.Millisecond,
	}
	manager := NewManager(config, nil)
	wrapper := NewImageTimeoutWrapper(manager)

	t.Run("operation completes in time", func(t *testing.T) {
		operation := func(ctx context.Context) error {
			time.Sleep(50 * time.Millisecond)
			return nil
		}

		wrappedOp := wrapper.WrapImageOperation(operation)
		err := wrappedOp(context.Background())
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("operation times out", func(t *testing.T) {
		operation := func(ctx context.Context) error {
			select {
			case <-time.After(200 * time.Millisecond):
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		wrappedOp := wrapper.WrapImageOperation(operation)
		err := wrappedOp(context.Background())
		if err == nil {
			t.Error("Expected timeout error")
		}
	})
}

func TestGracefulShutdown(t *testing.T) {
	config := &TimeoutConfig{
		ShutdownTimeout: 200 * time.Millisecond,
	}
	manager := NewManager(config, nil)
	shutdown := NewGracefulShutdown(manager)

	t.Run("all shutdown functions complete", func(t *testing.T) {
		var completed []int
		
		shutdown.AddShutdownFunc(func(ctx context.Context) error {
			time.Sleep(50 * time.Millisecond)
			completed = append(completed, 1)
			return nil
		})
		
		shutdown.AddShutdownFunc(func(ctx context.Context) error {
			time.Sleep(50 * time.Millisecond)
			completed = append(completed, 2)
			return nil
		})

		err := shutdown.Shutdown(context.Background())
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		if len(completed) != 2 {
			t.Errorf("Expected 2 functions to complete, got %d", len(completed))
		}
	})

	t.Run("shutdown times out", func(t *testing.T) {
		shutdownSlow := NewGracefulShutdown(manager)
		
		shutdownSlow.AddShutdownFunc(func(ctx context.Context) error {
			select {
			case <-time.After(300 * time.Millisecond): // Longer than shutdown timeout
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		})

		err := shutdownSlow.Shutdown(context.Background())
		if err == nil {
			t.Error("Expected timeout error")
		}
	})

	t.Run("shutdown function fails", func(t *testing.T) {
		shutdownFail := NewGracefulShutdown(manager)
		
		shutdownFail.AddShutdownFunc(func(ctx context.Context) error {
			return errors.New("shutdown failed")
		})
		
		shutdownFail.AddShutdownFunc(func(ctx context.Context) error {
			return nil
		})

		// Should complete despite one function failing
		err := shutdownFail.Shutdown(context.Background())
		if err != nil {
			t.Errorf("Expected no error even with failed function, got %v", err)
		}
	})
}

func TestCircuitBreaker(t *testing.T) {
	manager := NewManager(DefaultTimeoutConfig(), nil)
	cb := NewCircuitBreaker(manager, 3, 100*time.Millisecond)

	if cb.GetState() != StateClosed {
		t.Error("Circuit breaker should start closed")
	}

	t.Run("successful operations keep circuit closed", func(t *testing.T) {
		operation := func(ctx context.Context) error {
			return nil
		}

		for i := 0; i < 5; i++ {
			err := cb.Execute(context.Background(), operation)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if cb.GetState() != StateClosed {
				t.Error("Circuit should remain closed")
			}
		}
	})

	t.Run("failures open circuit", func(t *testing.T) {
		cbFail := NewCircuitBreaker(manager, 2, 100*time.Millisecond)
		operation := func(ctx context.Context) error {
			return errors.New("operation failed")
		}

		// First failure
		err := cbFail.Execute(context.Background(), operation)
		if err == nil {
			t.Error("Expected error")
		}
		if cbFail.GetState() != StateClosed {
			t.Error("Circuit should still be closed after first failure")
		}

		// Second failure should open circuit
		err = cbFail.Execute(context.Background(), operation)
		if err == nil {
			t.Error("Expected error")
		}
		if cbFail.GetState() != StateOpen {
			t.Error("Circuit should be open after max failures")
		}

		// Next operation should fail immediately
		err = cbFail.Execute(context.Background(), operation)
		if err == nil {
			t.Error("Expected circuit breaker error")
		}
	})

	t.Run("circuit recovers after timeout", func(t *testing.T) {
		cbRecover := NewCircuitBreaker(manager, 1, 50*time.Millisecond)
		
		// Cause failure to open circuit
		failOperation := func(ctx context.Context) error {
			return errors.New("operation failed")
		}
		cbRecover.Execute(context.Background(), failOperation)
		
		if cbRecover.GetState() != StateOpen {
			t.Error("Circuit should be open")
		}

		// Wait for timeout
		time.Sleep(100 * time.Millisecond)

		// Next operation should change state to half-open, then closed on success
		successOperation := func(ctx context.Context) error {
			return nil
		}
		err := cbRecover.Execute(context.Background(), successOperation)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		if cbRecover.GetState() != StateClosed {
			t.Error("Circuit should be closed after successful operation")
		}
	})
}

func TestCircuitBreakerGetStats(t *testing.T) {
	manager := NewManager(DefaultTimeoutConfig(), nil)
	cb := NewCircuitBreaker(manager, 3, 100*time.Millisecond)

	stats := cb.GetStats()
	
	if stats["state"] != "closed" {
		t.Error("Initial state should be closed")
	}
	if stats["failure_count"] != 0 {
		t.Error("Initial failure count should be 0")
	}
	if stats["max_failures"] != 3 {
		t.Error("Max failures should be 3")
	}
	if stats["timeout_ms"] != int64(100) {
		t.Error("Timeout should be 100ms")
	}
}

func TestCircuitStateString(t *testing.T) {
	tests := []struct {
		state    CircuitState
		expected string
	}{
		{StateClosed, "closed"},
		{StateOpen, "open"},
		{StateHalfOpen, "half-open"},
		{CircuitState(999), "unknown"},
	}

	for _, tt := range tests {
		if tt.state.String() != tt.expected {
			t.Errorf("Expected %s, got %s", tt.expected, tt.state.String())
		}
	}
}

func TestGetTimeoutForOperation(t *testing.T) {
	config := &TimeoutConfig{
		RequestTimeout:     30 * time.Second,
		DatabaseTimeout:    10 * time.Second,
		ImageTimeout:       15 * time.Second,
		ShutdownTimeout:    25 * time.Second,
		HealthCheckTimeout: 5 * time.Second,
	}
	manager := NewManager(config, nil)

	tests := []struct {
		operation string
		expected  time.Duration
	}{
		{"request", 30 * time.Second},
		{"database", 10 * time.Second},
		{"image_processing", 15 * time.Second},
		{"shutdown", 25 * time.Second},
		{"health_check", 5 * time.Second},
		{"unknown", 30 * time.Second}, // Should default to request timeout
	}

	for _, tt := range tests {
		result := manager.getTimeoutForOperation(tt.operation)
		if result != tt.expected {
			t.Errorf("For operation %s, expected %v, got %v", tt.operation, tt.expected, result)
		}
	}
}

func TestTimeoutWithLogging(t *testing.T) {
	logger := logging.New(logging.LevelDebug)
	manager := NewManager(DefaultTimeoutConfig(), logger)

	// Test that timeout creation and handling work with logger
	ctx, cancel := manager.WithRequestTimeout(context.Background())
	cancel()

	err := manager.HandleTimeout(ctx, "test", nil)
	if err != nil {
		t.Errorf("Expected no error with nil input, got %v", err)
	}
}

func BenchmarkTimeoutCreation(b *testing.B) {
	manager := NewManager(DefaultTimeoutConfig(), nil)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		timeoutCtx, cancel := manager.WithRequestTimeout(ctx)
		cancel()
		_ = timeoutCtx
	}
}

func BenchmarkCircuitBreakerExecution(b *testing.B) {
	manager := NewManager(DefaultTimeoutConfig(), nil)
	cb := NewCircuitBreaker(manager, 10, time.Second)
	
	operation := func(ctx context.Context) error {
		return nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cb.Execute(context.Background(), operation)
	}
}