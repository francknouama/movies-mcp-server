package support

import (
	"context"
)

// FaultInjector provides stubbed fault injection for BDD tests
// This is a simplified version that doesn't require testcontainers
type FaultInjector struct{}

// NewFaultInjector creates a new fault injector instance
func NewFaultInjector() *FaultInjector {
	return &FaultInjector{}
}

// ResourceMetrics represents resource usage metrics
type ResourceMetrics struct {
	MemoryUsageMB   int
	CPUUsagePercent float64
	ActiveConns     int
}

// ChaosConfig represents chaos engineering configuration
type ChaosConfig struct {
	PartialFailureRate float64
	DatabaseFailure    bool
	NetworkFailure     bool
	MemoryPressure     bool
}

// InjectDatabaseFailure simulates database failure (stubbed)
func (fi *FaultInjector) InjectDatabaseFailure(ctx context.Context) error {
	// Stubbed - does nothing in SQLite mode
	// In a real implementation, this would use chaos engineering tools
	return nil
}

// InjectMemoryPressure simulates memory pressure (stubbed)
func (fi *FaultInjector) InjectMemoryPressure(sizeMB int) error {
	// Stubbed - does nothing
	// In a real implementation, this would allocate memory to create pressure
	return nil
}

// InjectNetworkErrors simulates network errors (stubbed)
func (fi *FaultInjector) InjectNetworkErrors(targets []string) error {
	// Stubbed - does nothing
	// In a real implementation, this would use iptables or similar tools
	return nil
}

// InjectChaosConditions simulates various chaos conditions (stubbed)
func (fi *FaultInjector) InjectChaosConditions(ctx context.Context, config ChaosConfig) error {
	// Stubbed - does nothing
	// In a real implementation, this would coordinate multiple failure modes
	return nil
}

// MonitorResourceUsage monitors current resource usage
func (fi *FaultInjector) MonitorResourceUsage() (*ResourceMetrics, error) {
	// Return baseline metrics
	return &ResourceMetrics{
		MemoryUsageMB:   100,
		CPUUsagePercent: 10.0,
		ActiveConns:     5,
	}, nil
}

// Cleanup cleans up any injected faults
func (fi *FaultInjector) Cleanup() error {
	// Nothing to clean up in stubbed mode
	return nil
}
