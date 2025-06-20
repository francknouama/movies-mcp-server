package metrics

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"movies-mcp-server/pkg/logging"
)

// MetricType represents the type of metric
type MetricType string

const (
	CounterType   MetricType = "counter"
	GaugeType     MetricType = "gauge"
	HistogramType MetricType = "histogram"
	TimerType     MetricType = "timer"
)

// Metric represents a single metric
type Metric struct {
	Name      string            `json:"name"`
	Type      MetricType        `json:"type"`
	Value     float64           `json:"value"`
	Unit      string            `json:"unit"`
	Tags      map[string]string `json:"tags"`
	Timestamp time.Time         `json:"timestamp"`
}

// Metrics holds all application metrics
type Metrics struct {
	mu                    sync.RWMutex
	counters              map[string]*int64
	gauges                map[string]*float64
	histograms            map[string]*Histogram
	timers                map[string]*Timer
	logger                *logging.Logger
	startTime             time.Time
	lastReportTime        time.Time
	reportInterval        time.Duration
	
	// Built-in metrics
	RequestsTotal         *int64
	RequestsInFlight      *int64
	RequestDuration       *Histogram
	DatabaseOperations    *int64
	DatabaseErrors        *int64
	ImageOperations       *int64
	ImageProcessingTime   *Histogram
	MemoryUsage           *float64
	GoroutineCount        *float64
}

// Histogram tracks distribution of values
type Histogram struct {
	mu      sync.RWMutex
	buckets map[float64]int64
	count   int64
	sum     float64
}

// Timer tracks timing information
type Timer struct {
	histogram *Histogram
	startTime time.Time
}

// NewMetrics creates a new metrics instance
func NewMetrics(logger *logging.Logger, reportInterval time.Duration) *Metrics {
	m := &Metrics{
		counters:       make(map[string]*int64),
		gauges:         make(map[string]*float64),
		histograms:     make(map[string]*Histogram),
		timers:         make(map[string]*Timer),
		logger:         logger,
		startTime:      time.Now(),
		lastReportTime: time.Now(),
		reportInterval: reportInterval,
	}

	// Initialize built-in metrics
	m.RequestsTotal = m.NewCounter("requests_total", "Total number of requests")
	m.RequestsInFlight = m.NewCounter("requests_in_flight", "Number of requests currently being processed")
	m.RequestDuration = m.NewHistogram("request_duration_ms", "Request duration in milliseconds")
	m.DatabaseOperations = m.NewCounter("database_operations_total", "Total database operations")
	m.DatabaseErrors = m.NewCounter("database_errors_total", "Total database errors")
	m.ImageOperations = m.NewCounter("image_operations_total", "Total image operations")
	m.ImageProcessingTime = m.NewHistogram("image_processing_time_ms", "Image processing time in milliseconds")
	m.MemoryUsage = m.NewGauge("memory_usage_bytes", "Memory usage in bytes")
	m.GoroutineCount = m.NewGauge("goroutine_count", "Number of active goroutines")

	// Start background metric collection
	go m.collectSystemMetrics()
	go m.periodicReport()

	return m
}

// NewCounter creates a new counter metric
func (m *Metrics) NewCounter(name, description string) *int64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	counter := new(int64)
	m.counters[name] = counter
	return counter
}

// NewGauge creates a new gauge metric
func (m *Metrics) NewGauge(name, description string) *float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	gauge := new(float64)
	m.gauges[name] = gauge
	return gauge
}

// NewHistogram creates a new histogram metric
func (m *Metrics) NewHistogram(name, description string) *Histogram {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	histogram := &Histogram{
		buckets: make(map[float64]int64),
	}
	m.histograms[name] = histogram
	return histogram
}

// NewTimer creates a new timer metric
func (m *Metrics) NewTimer(name, description string) *Timer {
	histogram := m.NewHistogram(name+"_histogram", description+" histogram")
	return &Timer{
		histogram: histogram,
		startTime: time.Now(),
	}
}

// IncCounter increments a counter
func (m *Metrics) IncCounter(counter *int64) {
	atomic.AddInt64(counter, 1)
}

// AddCounter adds a value to a counter
func (m *Metrics) AddCounter(counter *int64, value int64) {
	atomic.AddInt64(counter, value)
}

// SetGauge sets a gauge value
func (m *Metrics) SetGauge(gauge *float64, value float64) {
	// Use atomic operations for float64 (requires some conversion)
	atomic.StoreUint64((*uint64)(unsafe.Pointer(gauge)), *(*uint64)(unsafe.Pointer(&value)))
}

// GetGauge gets a gauge value
func (m *Metrics) GetGauge(gauge *float64) float64 {
	bits := atomic.LoadUint64((*uint64)(unsafe.Pointer(gauge)))
	return *(*float64)(unsafe.Pointer(&bits))
}

// Observe adds a value to a histogram
func (h *Histogram) Observe(value float64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	h.count++
	h.sum += value
	
	// Define bucket boundaries (exponential buckets)
	buckets := []float64{1, 2, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000}
	
	for _, bucket := range buckets {
		if value <= bucket {
			h.buckets[bucket]++
		}
	}
}

// GetStats returns histogram statistics
func (h *Histogram) GetStats() (count int64, sum float64, buckets map[float64]int64) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	bucketsCopy := make(map[float64]int64)
	for k, v := range h.buckets {
		bucketsCopy[k] = v
	}
	
	return h.count, h.sum, bucketsCopy
}

// Stop stops a timer and records the duration
func (t *Timer) Stop() time.Duration {
	duration := time.Since(t.startTime)
	t.histogram.Observe(float64(duration.Milliseconds()))
	return duration
}

// StartRequestTimer starts a timer for a request
func (m *Metrics) StartRequestTimer() *Timer {
	m.IncCounter(m.RequestsTotal)
	m.IncCounter(m.RequestsInFlight)
	return &Timer{
		histogram: m.RequestDuration,
		startTime: time.Now(),
	}
}

// FinishRequestTimer finishes a request timer
func (m *Metrics) FinishRequestTimer(timer *Timer) {
	timer.Stop()
	atomic.AddInt64(m.RequestsInFlight, -1)
}

// RecordDatabaseOperation records a database operation
func (m *Metrics) RecordDatabaseOperation(duration time.Duration, err error) {
	m.IncCounter(m.DatabaseOperations)
	if err != nil {
		m.IncCounter(m.DatabaseErrors)
	}
}

// RecordImageOperation records an image operation
func (m *Metrics) RecordImageOperation(duration time.Duration, size int64) {
	m.IncCounter(m.ImageOperations)
	m.ImageProcessingTime.Observe(float64(duration.Milliseconds()))
}

// collectSystemMetrics collects system-level metrics
func (m *Metrics) collectSystemMetrics() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for range ticker.C {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		
		m.SetGauge(m.MemoryUsage, float64(memStats.Alloc))
		m.SetGauge(m.GoroutineCount, float64(runtime.NumGoroutine()))
	}
}

// periodicReport periodically reports metrics
func (m *Metrics) periodicReport() {
	if m.reportInterval == 0 {
		return // Reporting disabled
	}
	
	ticker := time.NewTicker(m.reportInterval)
	defer ticker.Stop()
	
	for range ticker.C {
		m.ReportMetrics()
	}
}

// ReportMetrics reports all current metrics
func (m *Metrics) ReportMetrics() {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	now := time.Now()
	
	// Report counters
	for name, counter := range m.counters {
		value := atomic.LoadInt64(counter)
		m.logger.LogPerformanceMetric(name, float64(value), "count", map[string]string{
			"type": "counter",
		})
	}
	
	// Report gauges
	for name, gauge := range m.gauges {
		value := m.GetGauge(gauge)
		m.logger.LogPerformanceMetric(name, value, "value", map[string]string{
			"type": "gauge",
		})
	}
	
	// Report histograms
	for name, histogram := range m.histograms {
		count, sum, buckets := histogram.GetStats()
		if count > 0 {
			avg := sum / float64(count)
			m.logger.LogPerformanceMetric(name+"_avg", avg, "ms", map[string]string{
				"type": "histogram",
				"stat": "average",
			})
			m.logger.LogPerformanceMetric(name+"_count", float64(count), "count", map[string]string{
				"type": "histogram",
				"stat": "count",
			})
			m.logger.LogPerformanceMetric(name+"_sum", sum, "ms", map[string]string{
				"type": "histogram",
				"stat": "sum",
			})
			
			// Report percentiles (simplified)
			for bucket, bucketCount := range buckets {
				if bucketCount > 0 {
					m.logger.LogPerformanceMetric(name+"_bucket", float64(bucketCount), "count", map[string]string{
						"type":   "histogram",
						"stat":   "bucket",
						"bucket": string(rune(bucket)),
					})
				}
			}
		}
	}
	
	m.lastReportTime = now
}

// GetSummary returns a summary of all metrics
func (m *Metrics) GetSummary() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	summary := make(map[string]interface{})
	
	// Add basic info
	summary["uptime_seconds"] = time.Since(m.startTime).Seconds()
	summary["last_report"] = m.lastReportTime.Format(time.RFC3339)
	
	// Add key metrics
	summary["requests_total"] = atomic.LoadInt64(m.RequestsTotal)
	summary["requests_in_flight"] = atomic.LoadInt64(m.RequestsInFlight)
	summary["database_operations_total"] = atomic.LoadInt64(m.DatabaseOperations)
	summary["database_errors_total"] = atomic.LoadInt64(m.DatabaseErrors)
	summary["image_operations_total"] = atomic.LoadInt64(m.ImageOperations)
	summary["memory_usage_bytes"] = m.GetGauge(m.MemoryUsage)
	summary["goroutine_count"] = m.GetGauge(m.GoroutineCount)
	
	// Add request duration stats
	if count, sum, _ := m.RequestDuration.GetStats(); count > 0 {
		summary["avg_request_duration_ms"] = sum / float64(count)
	}
	
	// Add image processing stats
	if count, sum, _ := m.ImageProcessingTime.GetStats(); count > 0 {
		summary["avg_image_processing_time_ms"] = sum / float64(count)
	}
	
	return summary
}

// Reset resets all metrics (useful for testing)
func (m *Metrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Reset all counters
	for _, counter := range m.counters {
		atomic.StoreInt64(counter, 0)
	}
	
	// Reset all gauges
	for _, gauge := range m.gauges {
		m.SetGauge(gauge, 0)
	}
	
	// Reset all histograms
	for _, histogram := range m.histograms {
		histogram.mu.Lock()
		histogram.count = 0
		histogram.sum = 0
		histogram.buckets = make(map[float64]int64)
		histogram.mu.Unlock()
	}
}

// MetricsMiddleware provides middleware for tracking request metrics
type MetricsMiddleware struct {
	metrics *Metrics
}

// NewMetricsMiddleware creates a new metrics middleware
func NewMetricsMiddleware(metrics *Metrics) *MetricsMiddleware {
	return &MetricsMiddleware{metrics: metrics}
}

// WrapHandler wraps a handler function with metrics tracking
func (mm *MetricsMiddleware) WrapHandler(handlerName string, handler func()) func() {
	return func() {
		timer := mm.metrics.StartRequestTimer()
		defer mm.metrics.FinishRequestTimer(timer)
		
		handler()
	}
}

// RequestContext holds request-specific metrics context
type RequestContext struct {
	RequestID string
	Method    string
	StartTime time.Time
	Metrics   *Metrics
}

// NewRequestContext creates a new request context
func NewRequestContext(requestID, method string, metrics *Metrics) *RequestContext {
	return &RequestContext{
		RequestID: requestID,
		Method:    method,
		StartTime: time.Now(),
		Metrics:   metrics,
	}
}

// Finish completes the request context and records metrics
func (rc *RequestContext) Finish(err error) {
	duration := time.Since(rc.StartTime)
	
	// Record request metrics
	rc.Metrics.RequestDuration.Observe(float64(duration.Milliseconds()))
	
	// Log performance metric
	tags := map[string]string{
		"method":     rc.Method,
		"request_id": rc.RequestID,
	}
	
	if err != nil {
		tags["status"] = "error"
	} else {
		tags["status"] = "success"
	}
	
	rc.Metrics.logger.LogPerformanceMetric(
		"request_duration",
		float64(duration.Milliseconds()),
		"ms",
		tags,
	)
}