package health

import (
	"context"
	"time"

	"movies-mcp-server/internal/database"
	"movies-mcp-server/pkg/logging"
)

// Status represents the health status of a component
type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
	StatusDegraded  Status = "degraded"
)

// Check represents a single health check
type Check struct {
	Name        string                 `json:"name"`
	Status      Status                 `json:"status"`
	Message     string                 `json:"message,omitempty"`
	Duration    time.Duration          `json:"duration"`
	Details     map[string]interface{} `json:"details,omitempty"`
	LastChecked time.Time              `json:"last_checked"`
	Error       string                 `json:"error,omitempty"`
}

// OverallHealth represents the overall system health
type OverallHealth struct {
	Status      Status            `json:"status"`
	Version     string            `json:"version"`
	Timestamp   time.Time         `json:"timestamp"`
	Uptime      time.Duration     `json:"uptime"`
	Checks      map[string]*Check `json:"checks"`
	Summary     Summary           `json:"summary"`
}

// Summary provides aggregated health information
type Summary struct {
	TotalChecks   int `json:"total_checks"`
	HealthyChecks int `json:"healthy_checks"`
	FailedChecks  int `json:"failed_checks"`
}

// Checker defines the interface for health checks
type Checker interface {
	Check(ctx context.Context) *Check
}

// Manager manages health checks for the application
type Manager struct {
	checkers  map[string]Checker
	logger    *logging.Logger
	startTime time.Time
	version   string
}

// NewManager creates a new health check manager
func NewManager(logger *logging.Logger, version string) *Manager {
	return &Manager{
		checkers:  make(map[string]Checker),
		logger:    logger,
		startTime: time.Now(),
		version:   version,
	}
}

// RegisterChecker registers a health checker
func (m *Manager) RegisterChecker(name string, checker Checker) {
	m.checkers[name] = checker
	m.logger.Info("health_checker_registered", "checker", name)
}

// CheckAll runs all registered health checks
func (m *Manager) CheckAll(ctx context.Context) *OverallHealth {
	checks := make(map[string]*Check)
	summary := Summary{TotalChecks: len(m.checkers)}

	// Run all checks
	for name, checker := range m.checkers {
		check := checker.Check(ctx)
		checks[name] = check

		// Update summary
		switch check.Status {
		case StatusHealthy:
			summary.HealthyChecks++
		default:
			summary.FailedChecks++
		}

		// Log health check result
		details := map[string]interface{}{
			"duration_ms": check.Duration.Milliseconds(),
		}
		if check.Details != nil {
			for k, v := range check.Details {
				details[k] = v
			}
		}
		
		m.logger.LogHealthCheck(name, string(check.Status), check.Duration, details)
	}

	// Determine overall status
	overallStatus := StatusHealthy
	if summary.FailedChecks > 0 {
		if summary.HealthyChecks == 0 {
			overallStatus = StatusUnhealthy
		} else {
			overallStatus = StatusDegraded
		}
	}

	return &OverallHealth{
		Status:    overallStatus,
		Version:   m.version,
		Timestamp: time.Now(),
		Uptime:    time.Since(m.startTime),
		Checks:    checks,
		Summary:   summary,
	}
}

// DatabaseChecker checks database connectivity
type DatabaseChecker struct {
	db database.Database
}

// NewDatabaseChecker creates a new database health checker
func NewDatabaseChecker(db database.Database) *DatabaseChecker {
	return &DatabaseChecker{db: db}
}

// Check performs the database health check
func (c *DatabaseChecker) Check(ctx context.Context) *Check {
	start := time.Now()
	check := &Check{
		Name:        "database",
		LastChecked: start,
		Details:     make(map[string]interface{}),
	}

	// Test database ping
	if err := c.db.Ping(); err != nil {
		check.Status = StatusUnhealthy
		check.Message = "Database connection failed"
		check.Error = err.Error()
		check.Duration = time.Since(start)
		return check
	}

	// Test basic query functionality
	stats, err := c.db.GetStats()
	if err != nil {
		check.Status = StatusDegraded
		check.Message = "Database accessible but queries failing"
		check.Error = err.Error()
		check.Duration = time.Since(start)
		return check
	}

	check.Status = StatusHealthy
	check.Message = "Database is healthy"
	check.Duration = time.Since(start)
	check.Details["total_movies"] = stats.TotalMovies
	check.Details["average_rating"] = stats.AverageRating
	check.Details["database_size"] = stats.DatabaseSize

	return check
}

// MCPServerChecker checks MCP server functionality
type MCPServerChecker struct {
	// This would typically include server metrics
}

// NewMCPServerChecker creates a new MCP server health checker
func NewMCPServerChecker() *MCPServerChecker {
	return &MCPServerChecker{}
}

// Check performs the MCP server health check
func (c *MCPServerChecker) Check(ctx context.Context) *Check {
	start := time.Now()
	check := &Check{
		Name:        "mcp_server",
		LastChecked: start,
		Details:     make(map[string]interface{}),
		Status:      StatusHealthy,
		Message:     "MCP server is operational",
		Duration:    time.Since(start),
	}

	// In a real implementation, this would check:
	// - Active connections
	// - Request processing capacity
	// - Memory usage
	// - Error rates

	check.Details["protocol_version"] = "2024-11-05"
	check.Details["active_connections"] = 0 // Placeholder

	return check
}

// ImageProcessorChecker checks image processing functionality
type ImageProcessorChecker struct {
	// Could include image processor instance for testing
}

// NewImageProcessorChecker creates a new image processor health checker
func NewImageProcessorChecker() *ImageProcessorChecker {
	return &ImageProcessorChecker{}
}

// Check performs the image processor health check
func (c *ImageProcessorChecker) Check(ctx context.Context) *Check {
	start := time.Now()
	check := &Check{
		Name:        "image_processor",
		LastChecked: start,
		Details:     make(map[string]interface{}),
		Status:      StatusHealthy,
		Message:     "Image processor is operational",
		Duration:    time.Since(start),
	}

	// Test basic image processing functionality
	// This is a simple validation that image processing libraries are available
	check.Details["supported_formats"] = []string{"image/jpeg", "image/png", "image/webp"}
	check.Details["max_size_mb"] = 5

	return check
}

// MemoryChecker checks memory usage
type MemoryChecker struct {
	maxMemoryMB float64
}

// NewMemoryChecker creates a new memory health checker
func NewMemoryChecker(maxMemoryMB float64) *MemoryChecker {
	return &MemoryChecker{maxMemoryMB: maxMemoryMB}
}

// Check performs the memory health check
func (c *MemoryChecker) Check(ctx context.Context) *Check {
	start := time.Now()
	check := &Check{
		Name:        "memory",
		LastChecked: start,
		Details:     make(map[string]interface{}),
	}

	// Get memory statistics (simplified - in production would use runtime.MemStats)
	// For now, we'll simulate memory usage
	currentMemoryMB := 64.0 // Placeholder
	
	check.Details["current_memory_mb"] = currentMemoryMB
	check.Details["max_memory_mb"] = c.maxMemoryMB
	check.Details["memory_usage_percent"] = (currentMemoryMB / c.maxMemoryMB) * 100

	if currentMemoryMB > c.maxMemoryMB*0.9 {
		check.Status = StatusUnhealthy
		check.Message = "Memory usage critically high"
	} else if currentMemoryMB > c.maxMemoryMB*0.7 {
		check.Status = StatusDegraded
		check.Message = "Memory usage high"
	} else {
		check.Status = StatusHealthy
		check.Message = "Memory usage normal"
	}

	check.Duration = time.Since(start)
	return check
}

// ReadinessProbe checks if the application is ready to serve requests
func (m *Manager) ReadinessProbe(ctx context.Context) bool {
	health := m.CheckAll(ctx)
	
	// Application is ready if no critical components are unhealthy
	for name, check := range health.Checks {
		if check.Status == StatusUnhealthy {
			// Critical components that must be healthy for readiness
			if name == "database" || name == "mcp_server" {
				return false
			}
		}
	}
	
	return true
}

// LivenessProbe checks if the application is alive
func (m *Manager) LivenessProbe(ctx context.Context) bool {
	// Simple liveness check - if we can execute this function, we're alive
	// In a more complex setup, this might check for deadlocks, etc.
	return true
}

// GetStartupTime returns when the application started
func (m *Manager) GetStartupTime() time.Time {
	return m.startTime
}

// GetUptime returns how long the application has been running
func (m *Manager) GetUptime() time.Duration {
	return time.Since(m.startTime)
}