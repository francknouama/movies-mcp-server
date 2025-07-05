package support

import (
	"context"
	"fmt"
	"net"
	"runtime"
	"sync"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

// FaultInjector provides mechanisms for injecting faults into the system for testing
type FaultInjector struct {
	dbContainer    testcontainers.Container
	networkBlocks  map[string]bool
	memoryPressure bool
	networkMutex   sync.RWMutex
	memoryMutex    sync.RWMutex
}

// NewFaultInjector creates a new fault injection utility
func NewFaultInjector() *FaultInjector {
	return &FaultInjector{
		networkBlocks: make(map[string]bool),
	}
}

// SetDatabaseContainer sets the database container for fault injection
func (fi *FaultInjector) SetDatabaseContainer(container testcontainers.Container) {
	fi.dbContainer = container
}

// InjectDatabaseFailure simulates database connection loss
func (fi *FaultInjector) InjectDatabaseFailure(ctx context.Context) error {
	if fi.dbContainer == nil {
		return fmt.Errorf("no database container available for fault injection")
	}

	// Stop the database container to simulate connection loss
	err := fi.dbContainer.Stop(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to stop database container: %w", err)
	}

	return nil
}

// RestoreDatabaseConnection restores database connectivity
func (fi *FaultInjector) RestoreDatabaseConnection(ctx context.Context) error {
	if fi.dbContainer == nil {
		return fmt.Errorf("no database container available")
	}

	// Start the database container
	err := fi.dbContainer.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed to start database container: %w", err)
	}

	// Wait for the database to be ready
	time.Sleep(5 * time.Second)
	return nil
}

// InjectNetworkErrors simulates network connectivity issues
func (fi *FaultInjector) InjectNetworkErrors(targetHosts []string) error {
	fi.networkMutex.Lock()
	defer fi.networkMutex.Unlock()

	for _, host := range targetHosts {
		fi.networkBlocks[host] = true
	}

	// In a real implementation, this would manipulate iptables or network interfaces
	// For testing purposes, we'll simulate by blocking specific connections
	return nil
}

// RestoreNetworkConnectivity restores network connectivity
func (fi *FaultInjector) RestoreNetworkConnectivity() {
	fi.networkMutex.Lock()
	defer fi.networkMutex.Unlock()

	fi.networkBlocks = make(map[string]bool)
}

// IsNetworkBlocked checks if a host is currently blocked
func (fi *FaultInjector) IsNetworkBlocked(host string) bool {
	fi.networkMutex.RLock()
	defer fi.networkMutex.RUnlock()

	return fi.networkBlocks[host]
}

// InjectMemoryPressure simulates memory pressure conditions
func (fi *FaultInjector) InjectMemoryPressure(targetMB int) error {
	fi.memoryMutex.Lock()
	defer fi.memoryMutex.Unlock()

	if fi.memoryPressure {
		return fmt.Errorf("memory pressure already active")
	}

	// Allocate large amounts of memory to simulate pressure
	go fi.createMemoryPressure(targetMB)
	fi.memoryPressure = true

	return nil
}

// ReleaseMemoryPressure releases memory pressure
func (fi *FaultInjector) ReleaseMemoryPressure() {
	fi.memoryMutex.Lock()
	defer fi.memoryMutex.Unlock()

	fi.memoryPressure = false
	runtime.GC() // Force garbage collection
}

// createMemoryPressure allocates memory to simulate pressure
func (fi *FaultInjector) createMemoryPressure(targetMB int) {
	// Allocate memory in chunks
	chunkSize := 1024 * 1024 // 1MB chunks
	chunks := make([][]byte, targetMB)

	for i := 0; i < targetMB && fi.memoryPressure; i++ {
		chunks[i] = make([]byte, chunkSize)
		// Fill with data to prevent optimization
		for j := range chunks[i] {
			chunks[i][j] = byte(j % 256)
		}
		time.Sleep(100 * time.Millisecond) // Gradual allocation
	}

	// Keep memory allocated while pressure is active
	for fi.memoryPressure {
		time.Sleep(1 * time.Second)
	}

	// Release memory and force garbage collection
	_ = chunks // Use the variable to avoid ineffectual assignment
	runtime.GC()
}

// InjectConnectionTimeout simulates connection timeouts
func (fi *FaultInjector) InjectConnectionTimeout(host string, port int, timeoutDuration time.Duration) error {
	// Create a connection that will timeout
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeoutDuration)
	if err != nil {
		return err // This is expected for timeout injection
	}
	defer conn.Close()

	return nil
}

// InjectSlowResponses simulates slow database or network responses
func (fi *FaultInjector) InjectSlowResponses(delayDuration time.Duration) {
	// Add artificial delay to simulate slow responses
	time.Sleep(delayDuration)
}

// InjectPartialFailures simulates partial system failures
func (fi *FaultInjector) InjectPartialFailures(failureRate float64) bool {
	// Simple probability-based failure injection
	// Returns true if operation should fail
	utilities := NewTestUtilities()
	return utilities.GenerateRandomFloat() < failureRate
}

// MonitorResourceUsage monitors system resource usage during fault injection
func (fi *FaultInjector) MonitorResourceUsage() (*ResourceMetrics, error) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Safe conversion from uint64 to int with bounds checking
	const maxInt = int(^uint(0) >> 1)

	// Convert to MB safely with explicit bounds checking
	var memoryUsedMB, memoryTotalMB, gcCount int

	memUsedBytes := m.Alloc / 1024 / 1024
	if memUsedBytes > uint64(maxInt) {
		memoryUsedMB = maxInt
	} else {
		memoryUsedMB = int(memUsedBytes)
	}

	memTotalBytes := m.Sys / 1024 / 1024
	if memTotalBytes > uint64(maxInt) {
		memoryTotalMB = maxInt
	} else {
		memoryTotalMB = int(memTotalBytes)
	}

	if m.NumGC > uint32(^uint32(0)>>1) {
		gcCount = int(^uint32(0) >> 1)
	} else {
		gcCount = int(m.NumGC)
	}

	metrics := &ResourceMetrics{
		MemoryUsedMB:   memoryUsedMB,
		MemoryTotalMB:  memoryTotalMB,
		GoroutineCount: runtime.NumGoroutine(),
		GCCount:        gcCount,
		Timestamp:      time.Now(),
	}

	return metrics, nil
}

// ResourceMetrics represents system resource metrics
type ResourceMetrics struct {
	MemoryUsedMB   int
	MemoryTotalMB  int
	GoroutineCount int
	GCCount        int
	Timestamp      time.Time
}

// InjectChaosConditions simulates multiple concurrent failure conditions
func (fi *FaultInjector) InjectChaosConditions(ctx context.Context, config ChaosConfig) error {
	var wg sync.WaitGroup
	errors := make(chan error, 4)

	if config.DatabaseFailure {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := fi.InjectDatabaseFailure(ctx); err != nil {
				errors <- fmt.Errorf("database failure injection failed: %w", err)
			}
		}()
	}

	if config.NetworkErrors && len(config.NetworkTargets) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := fi.InjectNetworkErrors(config.NetworkTargets); err != nil {
				errors <- fmt.Errorf("network error injection failed: %w", err)
			}
		}()
	}

	if config.MemoryPressure && config.MemoryPressureMB > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := fi.InjectMemoryPressure(config.MemoryPressureMB); err != nil {
				errors <- fmt.Errorf("memory pressure injection failed: %w", err)
			}
		}()
	}

	wg.Wait()
	close(errors)

	// Check for any errors
	for err := range errors {
		if err != nil {
			return err
		}
	}

	return nil
}

// RestoreAllSystems restores all systems from fault injection
func (fi *FaultInjector) RestoreAllSystems(ctx context.Context) error {
	var wg sync.WaitGroup
	errors := make(chan error, 3)

	// Restore database
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := fi.RestoreDatabaseConnection(ctx); err != nil {
			errors <- fmt.Errorf("database restoration failed: %w", err)
		}
	}()

	// Restore network
	wg.Add(1)
	go func() {
		defer wg.Done()
		fi.RestoreNetworkConnectivity()
	}()

	// Release memory pressure
	wg.Add(1)
	go func() {
		defer wg.Done()
		fi.ReleaseMemoryPressure()
	}()

	wg.Wait()
	close(errors)

	// Check for any errors
	for err := range errors {
		if err != nil {
			return err
		}
	}

	return nil
}

// ChaosConfig configures chaos engineering scenarios
type ChaosConfig struct {
	DatabaseFailure    bool
	NetworkErrors      bool
	NetworkTargets     []string
	MemoryPressure     bool
	MemoryPressureMB   int
	TimeoutDuration    time.Duration
	PartialFailureRate float64
}

// ValidationScenarios provides pre-configured fault injection scenarios
var ValidationScenarios = map[string]ChaosConfig{
	"database_failure": {
		DatabaseFailure: true,
	},
	"network_partition": {
		NetworkErrors:  true,
		NetworkTargets: []string{"localhost", "127.0.0.1"},
	},
	"memory_exhaustion": {
		MemoryPressure:   true,
		MemoryPressureMB: 100, // 100MB pressure
	},
	"partial_failure": {
		PartialFailureRate: 0.3, // 30% failure rate
	},
	"cascade_failure": {
		DatabaseFailure:    true,
		NetworkErrors:      true,
		NetworkTargets:     []string{"localhost"},
		MemoryPressure:     true,
		MemoryPressureMB:   50,
		PartialFailureRate: 0.5,
	},
}
