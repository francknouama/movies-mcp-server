package support

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"math"
	mathrand "math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

// AdvancedTestUtilities provides enhanced testing capabilities
type AdvancedTestUtilities struct {
	utilities  *TestUtilities
	metrics    *PerformanceMetrics
	monitoring *MonitoringService
	injector   *ErrorInjector
	generator  *DataGenerator
}

// PerformanceMetrics tracks detailed performance measurements
type PerformanceMetrics struct {
	mu              sync.RWMutex
	operations      map[string]*OperationMetrics
	memorySnapshots []MemorySnapshot
	resourceUsage   *ResourceUsage
	networkMetrics  *NetworkMetrics
}

// OperationMetrics tracks metrics for specific operations
type OperationMetrics struct {
	Name          string
	Count         int64
	TotalDuration time.Duration
	MinDuration   time.Duration
	MaxDuration   time.Duration
	ErrorCount    int64
	LastExecution time.Time
	Percentiles   map[int]time.Duration // P50, P95, P99
}

// MemorySnapshot captures memory state at a point in time
type MemorySnapshot struct {
	Timestamp  time.Time
	Alloc      uint64
	TotalAlloc uint64
	Sys        uint64
	NumGC      uint32
	HeapInuse  uint64
	StackInuse uint64
}

// ResourceUsage tracks system resource consumption
type ResourceUsage struct {
	CPUPercent    float64
	MemoryPercent float64
	DiskIOBytes   uint64
	NetworkBytes  uint64
	OpenFiles     int
	Goroutines    int
}

// NetworkMetrics tracks network-related performance
type NetworkMetrics struct {
	RequestCount    int64
	ResponseCount   int64
	BytesSent       int64
	BytesReceived   int64
	ConnectionCount int64
	ErrorCount      int64
}

// MonitoringService provides real-time monitoring capabilities
type MonitoringService struct {
	enabled  bool
	interval time.Duration
	stopCh   chan struct{}
	metrics  *PerformanceMetrics
	logger   *log.Logger
}

// ErrorInjector provides fault injection capabilities
type ErrorInjector struct {
	enabled       bool
	injectionRate float64 // 0.0 to 1.0
	errorTypes    []ErrorType
	activeErrors  map[string]error
	mu            sync.RWMutex
}

// ErrorType defines different types of errors that can be injected
type ErrorType int

const (
	ErrorTypeNetwork ErrorType = iota
	ErrorTypeDatabase
	ErrorTypeMemory
	ErrorTypeTimeout
	ErrorTypeConcurrency
	ErrorTypeValidation
)

// DataGenerator provides advanced test data generation
type DataGenerator struct {
	random        *mathrand.Rand
	templates     map[string]DataTemplate
	constraints   map[string]interface{}
	relationships map[string][]string
}

// DataTemplate defines a template for generating test data
type DataTemplate struct {
	Name       string
	Fields     map[string]FieldTemplate
	Count      int
	Variations []string
}

// FieldTemplate defines how to generate a specific field
type FieldTemplate struct {
	Type        string
	Min         interface{}
	Max         interface{}
	Options     []interface{}
	Pattern     string
	Constraints []string
}

// NewAdvancedTestUtilities creates a new instance
func NewAdvancedTestUtilities(utilities *TestUtilities) *AdvancedTestUtilities {
	return &AdvancedTestUtilities{
		utilities:  utilities,
		metrics:    NewPerformanceMetrics(),
		monitoring: NewMonitoringService(),
		injector:   NewErrorInjector(),
		generator:  NewDataGenerator(),
	}
}

// NewPerformanceMetrics creates a new performance metrics instance
func NewPerformanceMetrics() *PerformanceMetrics {
	return &PerformanceMetrics{
		operations:      make(map[string]*OperationMetrics),
		memorySnapshots: make([]MemorySnapshot, 0),
		resourceUsage:   &ResourceUsage{},
		networkMetrics:  &NetworkMetrics{},
	}
}

// NewMonitoringService creates a new monitoring service
func NewMonitoringService() *MonitoringService {
	return &MonitoringService{
		enabled:  false,
		interval: 1 * time.Second,
		stopCh:   make(chan struct{}),
		logger:   log.New(os.Stdout, "[MONITOR] ", log.LstdFlags),
	}
}

// NewErrorInjector creates a new error injector
func NewErrorInjector() *ErrorInjector {
	return &ErrorInjector{
		enabled:       false,
		injectionRate: 0.0,
		errorTypes:    []ErrorType{},
		activeErrors:  make(map[string]error),
	}
}

// NewDataGenerator creates a new data generator
func NewDataGenerator() *DataGenerator {
	return &DataGenerator{
		// #nosec G404 - weak random is acceptable for test data generation
		random:        mathrand.New(mathrand.NewSource(time.Now().UnixNano())),
		templates:     make(map[string]DataTemplate),
		constraints:   make(map[string]interface{}),
		relationships: make(map[string][]string),
	}
}

// Performance Monitoring Methods

// StartMonitoring begins performance monitoring
func (atu *AdvancedTestUtilities) StartMonitoring() error {
	atu.monitoring.enabled = true
	atu.monitoring.metrics = atu.metrics

	go atu.monitoring.monitoringLoop()

	return nil
}

// StopMonitoring stops performance monitoring
func (atu *AdvancedTestUtilities) StopMonitoring() error {
	if atu.monitoring.enabled {
		atu.monitoring.enabled = false
		close(atu.monitoring.stopCh)
	}

	return nil
}

// RecordOperation records metrics for an operation
func (atu *AdvancedTestUtilities) RecordOperation(name string, duration time.Duration, err error) {
	atu.metrics.mu.Lock()
	defer atu.metrics.mu.Unlock()

	if atu.metrics.operations[name] == nil {
		atu.metrics.operations[name] = &OperationMetrics{
			Name:        name,
			Percentiles: make(map[int]time.Duration),
		}
	}

	op := atu.metrics.operations[name]
	op.Count++
	op.TotalDuration += duration
	op.LastExecution = time.Now()

	if op.MinDuration == 0 || duration < op.MinDuration {
		op.MinDuration = duration
	}
	if duration > op.MaxDuration {
		op.MaxDuration = duration
	}

	if err != nil {
		op.ErrorCount++
	}
}

// TakeMemorySnapshot captures current memory state
func (atu *AdvancedTestUtilities) TakeMemorySnapshot() MemorySnapshot {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	snapshot := MemorySnapshot{
		Timestamp:  time.Now(),
		Alloc:      memStats.Alloc,
		TotalAlloc: memStats.TotalAlloc,
		Sys:        memStats.Sys,
		NumGC:      memStats.NumGC,
		HeapInuse:  memStats.HeapInuse,
		StackInuse: memStats.StackInuse,
	}

	atu.metrics.mu.Lock()
	atu.metrics.memorySnapshots = append(atu.metrics.memorySnapshots, snapshot)
	atu.metrics.mu.Unlock()

	return snapshot
}

// GetOperationMetrics returns metrics for a specific operation
func (atu *AdvancedTestUtilities) GetOperationMetrics(name string) *OperationMetrics {
	atu.metrics.mu.RLock()
	defer atu.metrics.mu.RUnlock()

	return atu.metrics.operations[name]
}

// GetMemoryGrowth calculates memory growth between two points
func (atu *AdvancedTestUtilities) GetMemoryGrowth(start, end MemorySnapshot) uint64 {
	return end.Alloc - start.Alloc
}

// Error Injection Methods

// EnableErrorInjection enables fault injection
func (atu *AdvancedTestUtilities) EnableErrorInjection(rate float64, errorTypes []ErrorType) {
	atu.injector.mu.Lock()
	defer atu.injector.mu.Unlock()

	atu.injector.enabled = true
	atu.injector.injectionRate = rate
	atu.injector.errorTypes = errorTypes
}

// DisableErrorInjection disables fault injection
func (atu *AdvancedTestUtilities) DisableErrorInjection() {
	atu.injector.mu.Lock()
	defer atu.injector.mu.Unlock()

	atu.injector.enabled = false
	atu.injector.activeErrors = make(map[string]error)
}

// ShouldInjectError determines if an error should be injected
func (atu *AdvancedTestUtilities) ShouldInjectError(operation string) (bool, error) {
	atu.injector.mu.RLock()
	defer atu.injector.mu.RUnlock()

	if !atu.injector.enabled {
		return false, nil
	}

	// #nosec G404 - weak random is acceptable for test error injection
	if mathrand.Float64() > atu.injector.injectionRate {
		return false, nil
	}

	// Check if there's an active error for this operation
	if err, exists := atu.injector.activeErrors[operation]; exists {
		return true, err
	}

	// Generate a new error
	if len(atu.injector.errorTypes) > 0 {
		// #nosec G404 - math/rand is acceptable for test utilities, not security-critical
		errorType := atu.injector.errorTypes[mathrand.Intn(len(atu.injector.errorTypes))]
		err := atu.generateError(errorType, operation)
		atu.injector.activeErrors[operation] = err
		return true, err
	}

	return false, nil
}

// generateError creates an error of the specified type
func (atu *AdvancedTestUtilities) generateError(errorType ErrorType, operation string) error {
	switch errorType {
	case ErrorTypeNetwork:
		return fmt.Errorf("network error during %s: connection timeout", operation)
	case ErrorTypeDatabase:
		return fmt.Errorf("database error during %s: connection lost", operation)
	case ErrorTypeMemory:
		return fmt.Errorf("memory error during %s: out of memory", operation)
	case ErrorTypeTimeout:
		return fmt.Errorf("timeout error during %s: operation exceeded time limit", operation)
	case ErrorTypeConcurrency:
		return fmt.Errorf("concurrency error during %s: resource lock timeout", operation)
	case ErrorTypeValidation:
		return fmt.Errorf("validation error during %s: invalid input data", operation)
	default:
		return fmt.Errorf("unknown error during %s", operation)
	}
}

// Data Generation Methods

// GenerateDataset creates a dataset based on a template
func (atu *AdvancedTestUtilities) GenerateDataset(templateName string, count int) ([]map[string]interface{}, error) {
	template, exists := atu.generator.templates[templateName]
	if !exists {
		return nil, fmt.Errorf("template %s not found", templateName)
	}

	dataset := make([]map[string]interface{}, count)

	for i := 0; i < count; i++ {
		item := make(map[string]interface{})

		for fieldName, fieldTemplate := range template.Fields {
			value, err := atu.generateFieldValue(fieldTemplate)
			if err != nil {
				return nil, fmt.Errorf("failed to generate field %s: %w", fieldName, err)
			}
			item[fieldName] = value
		}

		dataset[i] = item
	}

	return dataset, nil
}

// generateFieldValue generates a value based on field template
func (atu *AdvancedTestUtilities) generateFieldValue(template FieldTemplate) (interface{}, error) {
	switch template.Type {
	case "string":
		return atu.generateStringValue(template)
	case "integer":
		return atu.generateIntegerValue(template)
	case "float":
		return atu.generateFloatValue(template)
	case "boolean":
		return atu.generator.random.Intn(2) == 1, nil
	case "datetime":
		return atu.generateDateTimeValue(template)
	case "uuid":
		return atu.generateUUID()
	case "email":
		return atu.generateEmail()
	case "url":
		return atu.generateURL()
	default:
		return nil, fmt.Errorf("unsupported field type: %s", template.Type)
	}
}

// generateStringValue generates string values
func (atu *AdvancedTestUtilities) generateStringValue(template FieldTemplate) (string, error) {
	if len(template.Options) > 0 {
		option := template.Options[atu.generator.random.Intn(len(template.Options))]
		return fmt.Sprintf("%v", option), nil
	}

	minLen := 5
	maxLen := 20

	if template.Min != nil {
		if min, ok := template.Min.(int); ok {
			minLen = min
		}
	}
	if template.Max != nil {
		if max, ok := template.Max.(int); ok {
			maxLen = max
		}
	}

	length := minLen + atu.generator.random.Intn(maxLen-minLen+1)
	return atu.utilities.GenerateRandomString(length), nil
}

// generateIntegerValue generates integer values
func (atu *AdvancedTestUtilities) generateIntegerValue(template FieldTemplate) (int, error) {
	min := 0
	max := 1000

	if template.Min != nil {
		if minVal, ok := template.Min.(int); ok {
			min = minVal
		}
	}
	if template.Max != nil {
		if maxVal, ok := template.Max.(int); ok {
			max = maxVal
		}
	}

	return min + atu.generator.random.Intn(max-min+1), nil
}

// generateFloatValue generates float values
func (atu *AdvancedTestUtilities) generateFloatValue(template FieldTemplate) (float64, error) {
	min := 0.0
	max := 100.0

	if template.Min != nil {
		if minVal, ok := template.Min.(float64); ok {
			min = minVal
		}
	}
	if template.Max != nil {
		if maxVal, ok := template.Max.(float64); ok {
			max = maxVal
		}
	}

	return min + atu.generator.random.Float64()*(max-min), nil
}

// generateDateTimeValue generates datetime values
func (atu *AdvancedTestUtilities) generateDateTimeValue(template FieldTemplate) (time.Time, error) {
	now := time.Now()
	start := now.AddDate(-1, 0, 0) // 1 year ago
	end := now

	if template.Min != nil {
		if minTime, ok := template.Min.(time.Time); ok {
			start = minTime
		}
	}
	if template.Max != nil {
		if maxTime, ok := template.Max.(time.Time); ok {
			end = maxTime
		}
	}

	duration := end.Sub(start)
	randomDuration := time.Duration(atu.generator.random.Int63n(int64(duration)))

	return start.Add(randomDuration), nil
}

// generateUUID generates a UUID-like string
func (atu *AdvancedTestUtilities) generateUUID() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Set version and variant bits
	bytes[6] = (bytes[6] & 0x0f) | 0x40 // Version 4
	bytes[8] = (bytes[8] & 0x3f) | 0x80 // Variant 10

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:16]), nil
}

// generateEmail generates a random email address
func (atu *AdvancedTestUtilities) generateEmail() (string, error) {
	domains := []string{"test.com", "example.org", "demo.net", "sample.io"}
	username := atu.utilities.GenerateRandomString(8)
	domain := domains[atu.generator.random.Intn(len(domains))]

	return fmt.Sprintf("%s@%s", username, domain), nil
}

// generateURL generates a random URL
func (atu *AdvancedTestUtilities) generateURL() (string, error) {
	schemes := []string{"http", "https"}
	hosts := []string{"api.test.com", "service.example.org", "app.demo.net"}

	scheme := schemes[atu.generator.random.Intn(len(schemes))]
	host := hosts[atu.generator.random.Intn(len(hosts))]
	path := atu.utilities.GenerateRandomString(10)

	return fmt.Sprintf("%s://%s/%s", scheme, host, path), nil
}

// Utility Methods

// CreateDataTemplate creates a new data template
func (atu *AdvancedTestUtilities) CreateDataTemplate(name string, fields map[string]FieldTemplate) {
	atu.generator.templates[name] = DataTemplate{
		Name:   name,
		Fields: fields,
	}
}

// CreateLargeDataset creates a large dataset for performance testing
func (atu *AdvancedTestUtilities) CreateLargeDataset(size int, includeRelationships bool) ([]map[string]interface{}, error) {
	dataset := make([]map[string]interface{}, size)

	for i := 0; i < size; i++ {
		item := map[string]interface{}{
			"id":          i + 1,
			"title":       fmt.Sprintf("Large Dataset Movie %d", i+1),
			"director":    atu.utilities.GenerateRandomString(15),
			"year":        2000 + atu.generator.random.Intn(24),
			"rating":      5.0 + atu.generator.random.Float64()*5.0,
			"genre":       []string{"Action", "Drama", "Comedy", "Sci-Fi"}[atu.generator.random.Intn(4)],
			"description": atu.utilities.GenerateRandomString(200),
			"created_at":  time.Now().Add(-time.Duration(atu.generator.random.Intn(365*24)) * time.Hour),
		}

		if includeRelationships {
			// Add related actors
			actorCount := 1 + atu.generator.random.Intn(5)
			actors := make([]map[string]interface{}, actorCount)
			for j := 0; j < actorCount; j++ {
				actors[j] = map[string]interface{}{
					"id":   j + 1,
					"name": atu.utilities.GenerateRandomString(20),
				}
			}
			item["actors"] = actors
		}

		dataset[i] = item
	}

	return dataset, nil
}

// CreateConcurrentWorkload creates a workload for concurrent testing
func (atu *AdvancedTestUtilities) CreateConcurrentWorkload(operations []string, duration time.Duration) chan string {
	workloadCh := make(chan string, 1000)

	go func() {
		defer close(workloadCh)

		end := time.Now().Add(duration)
		for time.Now().Before(end) {
			operation := operations[atu.generator.random.Intn(len(operations))]

			select {
			case workloadCh <- operation:
			case <-time.After(100 * time.Millisecond):
				// Skip if channel is full
			}

			// Small delay between operations
			time.Sleep(time.Duration(atu.generator.random.Intn(100)) * time.Millisecond)
		}
	}()

	return workloadCh
}

// GenerateTestImage creates a base64-encoded test image
func (atu *AdvancedTestUtilities) GenerateTestImage(width, height int) (string, error) {
	// Create a simple test image (just random bytes for demonstration)
	imageSize := width * height * 3 // RGB
	imageData := make([]byte, imageSize)

	_, err := rand.Read(imageData)
	if err != nil {
		return "", fmt.Errorf("failed to generate image data: %w", err)
	}

	return base64.StdEncoding.EncodeToString(imageData), nil
}

// CalculatePercentile calculates the specified percentile from a list of durations
func (atu *AdvancedTestUtilities) CalculatePercentile(durations []time.Duration, percentile int) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	// Sort durations (simple bubble sort for small datasets)
	for i := 0; i < len(durations)-1; i++ {
		for j := 0; j < len(durations)-i-1; j++ {
			if durations[j] > durations[j+1] {
				durations[j], durations[j+1] = durations[j+1], durations[j]
			}
		}
	}

	index := int(math.Ceil(float64(percentile)/100.0*float64(len(durations)))) - 1
	if index < 0 {
		index = 0
	}
	if index >= len(durations) {
		index = len(durations) - 1
	}

	return durations[index]
}

// monitoring loop for the monitoring service
func (ms *MonitoringService) monitoringLoop() {
	ticker := time.NewTicker(ms.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if !ms.enabled {
				return
			}

			// Collect current metrics
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)

			snapshot := MemorySnapshot{
				Timestamp:  time.Now(),
				Alloc:      memStats.Alloc,
				TotalAlloc: memStats.TotalAlloc,
				Sys:        memStats.Sys,
				NumGC:      memStats.NumGC,
				HeapInuse:  memStats.HeapInuse,
				StackInuse: memStats.StackInuse,
			}

			ms.metrics.mu.Lock()
			ms.metrics.memorySnapshots = append(ms.metrics.memorySnapshots, snapshot)

			// Keep only last 1000 snapshots to prevent memory growth
			if len(ms.metrics.memorySnapshots) > 1000 {
				ms.metrics.memorySnapshots = ms.metrics.memorySnapshots[1:]
			}
			ms.metrics.mu.Unlock()

		case <-ms.stopCh:
			return
		}
	}
}
