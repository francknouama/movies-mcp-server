package timeout

import (
	"context"
	"fmt"
	"time"

	"movies-mcp-server/pkg/errors"
	"movies-mcp-server/pkg/logging"
)

// TimeoutConfig holds timeout configuration
type TimeoutConfig struct {
	// RequestTimeout is the maximum time allowed for a single request
	RequestTimeout time.Duration
	// DatabaseTimeout is the maximum time allowed for database operations
	DatabaseTimeout time.Duration
	// ImageTimeout is the maximum time allowed for image processing
	ImageTimeout time.Duration
	// ShutdownTimeout is the maximum time allowed for graceful shutdown
	ShutdownTimeout time.Duration
	// HealthCheckTimeout is the maximum time allowed for health checks
	HealthCheckTimeout time.Duration
}

// DefaultTimeoutConfig returns default timeout configuration
func DefaultTimeoutConfig() *TimeoutConfig {
	return &TimeoutConfig{
		RequestTimeout:     30 * time.Second,
		DatabaseTimeout:    10 * time.Second,
		ImageTimeout:       15 * time.Second,
		ShutdownTimeout:    30 * time.Second,
		HealthCheckTimeout: 5 * time.Second,
	}
}

// Manager manages timeouts for various operations
type Manager struct {
	config *TimeoutConfig
	logger *logging.Logger
}

// NewManager creates a new timeout manager
func NewManager(config *TimeoutConfig, logger *logging.Logger) *Manager {
	if config == nil {
		config = DefaultTimeoutConfig()
	}
	return &Manager{
		config: config,
		logger: logger,
	}
}

// WithRequestTimeout creates a context with request timeout
func (m *Manager) WithRequestTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	return m.withTimeout(parent, m.config.RequestTimeout, "request")
}

// WithDatabaseTimeout creates a context with database timeout
func (m *Manager) WithDatabaseTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	return m.withTimeout(parent, m.config.DatabaseTimeout, "database")
}

// WithImageTimeout creates a context with image processing timeout
func (m *Manager) WithImageTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	return m.withTimeout(parent, m.config.ImageTimeout, "image_processing")
}

// WithHealthCheckTimeout creates a context with health check timeout
func (m *Manager) WithHealthCheckTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	return m.withTimeout(parent, m.config.HealthCheckTimeout, "health_check")
}

// WithShutdownTimeout creates a context with shutdown timeout
func (m *Manager) WithShutdownTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	return m.withTimeout(parent, m.config.ShutdownTimeout, "shutdown")
}

// WithCustomTimeout creates a context with custom timeout
func (m *Manager) WithCustomTimeout(parent context.Context, timeout time.Duration, operation string) (context.Context, context.CancelFunc) {
	return m.withTimeout(parent, timeout, operation)
}

// withTimeout is the internal method to create timeouts
func (m *Manager) withTimeout(parent context.Context, timeout time.Duration, operation string) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(parent, timeout)
	
	// Log timeout creation if logger is available
	if m.logger != nil {
		m.logger.Debug("timeout_created",
			"operation", operation,
			"timeout_ms", timeout.Milliseconds(),
		)
	}

	return ctx, cancel
}

// HandleTimeout handles context timeout errors and creates appropriate application errors
func (m *Manager) HandleTimeout(ctx context.Context, operation string, err error) error {
	if err == nil {
		return nil
	}

	// Check if the error is due to context timeout
	if ctx.Err() == context.DeadlineExceeded {
		timeout := m.getTimeoutForOperation(operation)
		
		if m.logger != nil {
			m.logger.Warn("operation_timeout",
				"operation", operation,
				"timeout_ms", timeout.Milliseconds(),
			)
		}

		return errors.NewTimeoutError(operation, timeout).
			WithSeverity(errors.SeverityHigh).
			WithComponent("timeout_manager")
	}

	// Check if the error is due to context cancellation
	if ctx.Err() == context.Canceled {
		if m.logger != nil {
			m.logger.Debug("operation_cancelled", "operation", operation)
		}
		return errors.NewApplicationError(errors.InternalError, "Operation cancelled").
			WithSeverity(errors.SeverityMedium).
			WithComponent("timeout_manager").
			WithDetails(map[string]interface{}{
				"operation": operation,
				"reason":    "cancelled",
			})
	}

	// Return the original error if not timeout/cancellation related
	return err
}

// getTimeoutForOperation returns the appropriate timeout for an operation
func (m *Manager) getTimeoutForOperation(operation string) time.Duration {
	switch operation {
	case "request":
		return m.config.RequestTimeout
	case "database":
		return m.config.DatabaseTimeout
	case "image_processing":
		return m.config.ImageTimeout
	case "health_check":
		return m.config.HealthCheckTimeout
	case "shutdown":
		return m.config.ShutdownTimeout
	default:
		return m.config.RequestTimeout // Default fallback
	}
}

// RequestTimeoutMiddleware provides timeout middleware for requests
type RequestTimeoutMiddleware struct {
	manager *Manager
}

// NewRequestTimeoutMiddleware creates a new request timeout middleware
func NewRequestTimeoutMiddleware(manager *Manager) *RequestTimeoutMiddleware {
	return &RequestTimeoutMiddleware{manager: manager}
}

// WrapHandler wraps a handler function with timeout handling
func (rtm *RequestTimeoutMiddleware) WrapHandler(handler func(ctx context.Context) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		timeoutCtx, cancel := rtm.manager.WithRequestTimeout(ctx)
		defer cancel()

		// Create a channel to receive the result
		resultChan := make(chan error, 1)

		// Run the handler in a goroutine
		go func() {
			defer func() {
				if r := recover(); r != nil {
					resultChan <- fmt.Errorf("handler panicked: %v", r)
				}
			}()
			resultChan <- handler(timeoutCtx)
		}()

		// Wait for either completion or timeout
		select {
		case err := <-resultChan:
			return rtm.manager.HandleTimeout(timeoutCtx, "request", err)
		case <-timeoutCtx.Done():
			return rtm.manager.HandleTimeout(timeoutCtx, "request", timeoutCtx.Err())
		}
	}
}

// DatabaseTimeoutWrapper provides timeout handling for database operations
type DatabaseTimeoutWrapper struct {
	manager *Manager
}

// NewDatabaseTimeoutWrapper creates a new database timeout wrapper
func NewDatabaseTimeoutWrapper(manager *Manager) *DatabaseTimeoutWrapper {
	return &DatabaseTimeoutWrapper{manager: manager}
}

// WrapDatabaseOperation wraps a database operation with timeout handling
func (dtw *DatabaseTimeoutWrapper) WrapDatabaseOperation(operation func(ctx context.Context) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		timeoutCtx, cancel := dtw.manager.WithDatabaseTimeout(ctx)
		defer cancel()

		// Execute the operation
		err := operation(timeoutCtx)
		return dtw.manager.HandleTimeout(timeoutCtx, "database", err)
	}
}

// ImageTimeoutWrapper provides timeout handling for image operations
type ImageTimeoutWrapper struct {
	manager *Manager
}

// NewImageTimeoutWrapper creates a new image timeout wrapper
func NewImageTimeoutWrapper(manager *Manager) *ImageTimeoutWrapper {
	return &ImageTimeoutWrapper{manager: manager}
}

// WrapImageOperation wraps an image operation with timeout handling
func (itw *ImageTimeoutWrapper) WrapImageOperation(operation func(ctx context.Context) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		timeoutCtx, cancel := itw.manager.WithImageTimeout(ctx)
		defer cancel()

		// Execute the operation
		err := operation(timeoutCtx)
		return itw.manager.HandleTimeout(timeoutCtx, "image_processing", err)
	}
}

// GracefulShutdown handles graceful shutdown with timeout
type GracefulShutdown struct {
	manager       *Manager
	shutdownFuncs []func(ctx context.Context) error
}

// NewGracefulShutdown creates a new graceful shutdown manager
func NewGracefulShutdown(manager *Manager) *GracefulShutdown {
	return &GracefulShutdown{
		manager:       manager,
		shutdownFuncs: make([]func(ctx context.Context) error, 0),
	}
}

// AddShutdownFunc adds a function to be called during shutdown
func (gs *GracefulShutdown) AddShutdownFunc(fn func(ctx context.Context) error) {
	gs.shutdownFuncs = append(gs.shutdownFuncs, fn)
}

// Shutdown performs graceful shutdown with timeout
func (gs *GracefulShutdown) Shutdown(ctx context.Context) error {
	shutdownCtx, cancel := gs.manager.WithShutdownTimeout(ctx)
	defer cancel()

	if gs.manager.logger != nil {
		gs.manager.logger.Info("graceful_shutdown_started",
			"timeout_ms", gs.manager.config.ShutdownTimeout.Milliseconds(),
			"shutdown_funcs_count", len(gs.shutdownFuncs),
		)
	}

	// Execute all shutdown functions
	for i, fn := range gs.shutdownFuncs {
		// Check if context is already done before starting function
		select {
		case <-shutdownCtx.Done():
			if gs.manager.logger != nil {
				gs.manager.logger.Warn("shutdown_timeout",
					"completed_funcs", i,
					"total_funcs", len(gs.shutdownFuncs),
				)
			}
			return gs.manager.HandleTimeout(shutdownCtx, "shutdown", shutdownCtx.Err())
		default:
		}

		// Execute shutdown function
		if err := fn(shutdownCtx); err != nil {
			if gs.manager.logger != nil {
				gs.manager.logger.Error("shutdown_func_failed",
					"func_index", i,
					"error", err.Error(),
				)
			}
			// Continue with other shutdown functions even if one fails
		}

		// Check if we're out of time after executing the function
		select {
		case <-shutdownCtx.Done():
			if gs.manager.logger != nil {
				gs.manager.logger.Warn("shutdown_timeout",
					"completed_funcs", i+1,
					"total_funcs", len(gs.shutdownFuncs),
				)
			}
			return gs.manager.HandleTimeout(shutdownCtx, "shutdown", shutdownCtx.Err())
		default:
		}
	}

	if gs.manager.logger != nil {
		gs.manager.logger.Info("graceful_shutdown_completed")
	}

	return nil
}

// CircuitBreaker provides circuit breaker functionality with timeout support
type CircuitBreaker struct {
	manager         *Manager
	failureCount    int
	maxFailures     int
	timeout         time.Duration
	lastFailureTime time.Time
	state           CircuitState
}

// CircuitState represents the state of a circuit breaker
type CircuitState int

const (
	StateClosed CircuitState = iota
	StateOpen
	StateHalfOpen
)

// String returns the string representation of circuit state
func (s CircuitState) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(manager *Manager, maxFailures int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		manager:     manager,
		maxFailures: maxFailures,
		timeout:     timeout,
		state:       StateClosed,
	}
}

// Execute executes an operation with circuit breaker protection
func (cb *CircuitBreaker) Execute(ctx context.Context, operation func(ctx context.Context) error) error {
	// Check circuit state
	if cb.state == StateOpen {
		if time.Since(cb.lastFailureTime) > cb.timeout {
			cb.state = StateHalfOpen
			if cb.manager.logger != nil {
				cb.manager.logger.Info("circuit_breaker_half_open")
			}
		} else {
			return errors.NewApplicationError(errors.ServiceUnavailable, "Circuit breaker is open").
				WithSeverity(errors.SeverityHigh).
				WithComponent("circuit_breaker").
				WithDetails(map[string]interface{}{
					"state":            cb.state.String(),
					"failure_count":    cb.failureCount,
					"max_failures":     cb.maxFailures,
					"last_failure_ago": time.Since(cb.lastFailureTime).String(),
				})
		}
	}

	// Execute the operation
	err := operation(ctx)

	// Handle result
	if err != nil {
		cb.recordFailure()
		return err
	}

	cb.recordSuccess()
	return nil
}

// recordFailure records a failure and potentially opens the circuit
func (cb *CircuitBreaker) recordFailure() {
	cb.failureCount++
	cb.lastFailureTime = time.Now()

	if cb.failureCount >= cb.maxFailures {
		cb.state = StateOpen
		if cb.manager.logger != nil {
			cb.manager.logger.Warn("circuit_breaker_opened",
				"failure_count", cb.failureCount,
				"max_failures", cb.maxFailures,
			)
		}
	}
}

// recordSuccess records a success and potentially closes the circuit
func (cb *CircuitBreaker) recordSuccess() {
	cb.failureCount = 0
	if cb.state == StateHalfOpen {
		cb.state = StateClosed
		if cb.manager.logger != nil {
			cb.manager.logger.Info("circuit_breaker_closed")
		}
	}
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() CircuitState {
	return cb.state
}

// GetStats returns circuit breaker statistics
func (cb *CircuitBreaker) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"state":            cb.state.String(),
		"failure_count":    cb.failureCount,
		"max_failures":     cb.maxFailures,
		"last_failure_ago": time.Since(cb.lastFailureTime).String(),
		"timeout_ms":       cb.timeout.Milliseconds(),
	}
}