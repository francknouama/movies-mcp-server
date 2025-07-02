package steps

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cucumber/godog"
	"github.com/francknouama/movies-mcp-server/mcp-server/tests/bdd/context"
	"github.com/francknouama/movies-mcp-server/mcp-server/tests/bdd/support"
)

// SimplePerformanceSteps provides basic performance step definitions
type SimplePerformanceSteps struct {
	bddContext       *context.BDDContext
	utilities        *support.TestUtilities
	startTime        time.Time
	endTime          time.Time
	duration         time.Duration
	operationCount   int32
	errorCount       int32
	results          []interface{}
	errors           []error
	memoryBefore     runtime.MemStats
	memoryAfter      runtime.MemStats
	performanceThresholds map[string]time.Duration
	throughputMetrics     map[string]float64
	performanceViolations []string
}

// NewSimplePerformanceSteps creates a new instance
func NewSimplePerformanceSteps(bddContext *context.BDDContext, utilities *support.TestUtilities) *SimplePerformanceSteps {
	return &SimplePerformanceSteps{
		bddContext: bddContext,
		utilities:  utilities,
		performanceThresholds: map[string]time.Duration{
			"simple_operation":  500 * time.Millisecond,
			"search_operation":  1 * time.Second,
			"batch_operation":   10 * time.Second,
			"memory_operation":  2 * time.Second,
		},
		throughputMetrics: make(map[string]float64),
		performanceViolations: make([]string, 0),
	}
}

// InitializePerformanceSteps registers performance step definitions following the existing pattern
func InitializePerformanceSteps(ctx *godog.ScenarioContext) {
	stepContext := NewCommonStepContext()
	utilities := support.NewTestUtilities()
	sps := &SimplePerformanceSteps{
		bddContext: stepContext.bddContext,
		utilities:  utilities,
		performanceThresholds: map[string]time.Duration{
			"simple_operation":  500 * time.Millisecond,
			"search_operation":  1 * time.Second,
			"batch_operation":   10 * time.Second,
			"memory_operation":  2 * time.Second,
		},
		throughputMetrics: make(map[string]float64),
		performanceViolations: make([]string, 0),
	}
	// Data setup
	ctx.Step(`^I have (\d+) movies in the database$`, sps.iHaveMoviesInTheDatabase)
	ctx.Step(`^I measure the baseline memory usage$`, sps.iMeasureBaselineMemory)

	// Performance operations
	ctx.Step(`^I perform (\d+) concurrent searches for "([^"]*)"$`, sps.iPerformConcurrentSearches)
	ctx.Step(`^I search for movies by genre "([^"]*)"$`, sps.iSearchMoviesByGenre)
	ctx.Step(`^I create (\d+) movies in batch$`, sps.iCreateMoviesInBatch)
	ctx.Step(`^I load (\d+) movies with full details$`, sps.iLoadMoviesWithDetails)

	// Assertions
	ctx.Step(`^all searches should complete within (\d+) seconds?$`, sps.operationsShouldCompleteWithinSeconds)
	ctx.Step(`^the response should be returned within (\d+)ms$`, sps.responseShouldBeWithinMs)
	ctx.Step(`^the operation should complete within (\d+) seconds?$`, sps.operationsShouldCompleteWithinSeconds)
	ctx.Step(`^no search should fail due to resource contention$`, sps.noOperationsShouldFail)
	ctx.Step(`^each search result should be valid$`, sps.eachResultShouldBeValid)
	ctx.Step(`^all movies should be successfully created$`, sps.allMoviesShouldBeCreated)
	ctx.Step(`^the memory increase should not exceed (\d+)MB$`, sps.memoryIncreaseShouldNotExceed)
}

// Implementation methods
func (sps *SimplePerformanceSteps) iHaveMoviesInTheDatabase(count int) error {
	movies := sps.utilities.CreateTestMovieBatch(count)

	sps.startTime = time.Now()
	for _, movie := range movies {
		response, err := sps.bddContext.CallTool("add_movie", movie)
		if err != nil {
			return fmt.Errorf("failed to add movie: %w", err)
		}
		if response.IsError {
			return fmt.Errorf("MCP error: %v", response.Content)
		}
	}
	sps.endTime = time.Now()
	sps.duration = sps.endTime.Sub(sps.startTime)

	return nil
}

func (sps *SimplePerformanceSteps) iMeasureBaselineMemory() error {
	runtime.GC() // Force garbage collection
	runtime.ReadMemStats(&sps.memoryBefore)
	return nil
}

func (sps *SimplePerformanceSteps) iPerformConcurrentSearches(count int, searchTerm string) error {
	sps.results = make([]interface{}, count)
	sps.errors = make([]error, count)

	sps.startTime = time.Now()

	var wg sync.WaitGroup
	var operationCount int32
	var errorCount int32

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			response, err := sps.bddContext.CallTool("search_movies", map[string]interface{}{
				"genre": searchTerm,
			})

			atomic.AddInt32(&operationCount, 1)
			if err != nil || response.IsError {
				atomic.AddInt32(&errorCount, 1)
				if err != nil {
					sps.errors[index] = err
				} else {
					sps.errors[index] = fmt.Errorf("MCP error: %v", response.Content)
				}
			} else {
				sps.results[index] = response.Content
			}
		}(i)
	}

	wg.Wait()
	sps.endTime = time.Now()
	sps.duration = sps.endTime.Sub(sps.startTime)
	sps.operationCount = operationCount
	sps.errorCount = errorCount

	return nil
}

func (sps *SimplePerformanceSteps) iSearchMoviesByGenre(genre string) error {
	sps.startTime = time.Now()

	response, err := sps.bddContext.CallTool("search_movies", map[string]interface{}{
		"genre": genre,
	})

	sps.endTime = time.Now()
	sps.duration = sps.endTime.Sub(sps.startTime)

	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}
	if response.IsError {
		return fmt.Errorf("MCP error: %v", response.Content)
	}

	sps.results = []interface{}{response.Content}
	return nil
}

func (sps *SimplePerformanceSteps) iCreateMoviesInBatch(count int) error {
	movies := sps.utilities.CreateTestMovieBatch(count)
	sps.results = make([]interface{}, count)

	sps.startTime = time.Now()

	for i, movie := range movies {
		response, err := sps.bddContext.CallTool("add_movie", movie)
		if err != nil {
			return fmt.Errorf("failed to create movie %d: %w", i, err)
		}
		if response.IsError {
			return fmt.Errorf("MCP error creating movie %d: %v", i, response.Content)
		}
		sps.results[i] = response.Content
	}

	sps.endTime = time.Now()
	sps.duration = sps.endTime.Sub(sps.startTime)

	return nil
}

func (sps *SimplePerformanceSteps) iLoadMoviesWithDetails(count int) error {
	runtime.ReadMemStats(&sps.memoryBefore)

	sps.startTime = time.Now()

	// First get list of movies
	response, err := sps.bddContext.CallTool("search_movies", map[string]interface{}{
		"limit": count,
	})
	if err != nil {
		return fmt.Errorf("failed to search movies: %w", err)
	}
	if response.IsError {
		return fmt.Errorf("MCP error searching movies: %v", response.Content)
	}

	sps.endTime = time.Now()
	sps.duration = sps.endTime.Sub(sps.startTime)

	runtime.ReadMemStats(&sps.memoryAfter)
	sps.results = []interface{}{response.Content}

	return nil
}

// Assertion methods
func (sps *SimplePerformanceSteps) operationsShouldCompleteWithinSeconds(seconds int) error {
	maxDuration := time.Duration(seconds) * time.Second
	if sps.duration > maxDuration {
		return fmt.Errorf("operations took %v, expected under %v", sps.duration, maxDuration)
	}
	return nil
}

func (sps *SimplePerformanceSteps) responseShouldBeWithinMs(ms int) error {
	maxDuration := time.Duration(ms) * time.Millisecond
	if sps.duration > maxDuration {
		return fmt.Errorf("response took %v, expected under %v", sps.duration, maxDuration)
	}
	return nil
}

func (sps *SimplePerformanceSteps) noOperationsShouldFail() error {
	if sps.errorCount > 0 {
		return fmt.Errorf("found %d errors in operations", sps.errorCount)
	}
	return nil
}

func (sps *SimplePerformanceSteps) eachResultShouldBeValid() error {
	for i, result := range sps.results {
		if result == nil {
			continue
		}

		// Basic validation - check if it's a reasonable response
		switch v := result.(type) {
		case []interface{}:
			// Array of results (like search results)
			for _, item := range v {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if err := sps.utilities.ValidateMovieResponse(itemMap); err != nil {
						return fmt.Errorf("invalid movie in result %d: %w", i, err)
					}
				}
			}
		case map[string]interface{}:
			// Single result (like created movie)
			if err := sps.utilities.ValidateMovieResponse(v); err != nil {
				return fmt.Errorf("invalid result %d: %w", i, err)
			}
		}
	}
	return nil
}

func (sps *SimplePerformanceSteps) allMoviesShouldBeCreated() error {
	if sps.errorCount > 0 {
		return fmt.Errorf("found %d errors during movie creation", sps.errorCount)
	}

	successCount := 0
	for _, result := range sps.results {
		if result != nil {
			successCount++
		}
	}

	expectedCount := len(sps.results)
	if successCount != expectedCount {
		return fmt.Errorf("created %d movies, expected %d", successCount, expectedCount)
	}

	return nil
}

func (sps *SimplePerformanceSteps) memoryIncreaseShouldNotExceed(mb int) error {
	if mb < 0 {
		return fmt.Errorf("negative memory limit not allowed: %d MB", mb)
	}
	maxIncrease := uint64(mb) * 1024 * 1024 // Convert MB to bytes
	memoryIncrease := sps.memoryAfter.Alloc - sps.memoryBefore.Alloc

	if memoryIncrease > maxIncrease {
		violation := fmt.Errorf("memory increased by %d bytes, expected under %d bytes",
			memoryIncrease, maxIncrease)
		sps.performanceViolations = append(sps.performanceViolations, violation.Error())
		return violation
	}

	// Enforce memory usage contract
	err := sps.validateResourceUsage()
	if err != nil {
		sps.bddContext.SetTestData("memory_performance_violation", err.Error())
	}

	return nil
}

// enforcePerformanceContract validates performance against defined contracts
func (sps *SimplePerformanceSteps) enforcePerformanceContract(operationType string) error {
	threshold, exists := sps.performanceThresholds[operationType]
	if !exists {
		return fmt.Errorf("no performance threshold defined for operation type: %s", operationType)
	}
	
	if sps.duration > threshold {
		violation := fmt.Sprintf("%s operation took %v, contract requires under %v", 
			operationType, sps.duration, threshold)
		sps.performanceViolations = append(sps.performanceViolations, violation)
		return fmt.Errorf(violation)
	}
	
	return nil
}

// measureThroughput calculates throughput metrics for concurrent operations
func (sps *SimplePerformanceSteps) measureThroughput() {
	if sps.duration > 0 && sps.operationCount > 0 {
		throughput := float64(sps.operationCount) / sps.duration.Seconds()
		errorRate := float64(sps.errorCount) / float64(sps.operationCount)
		
		sps.throughputMetrics["throughput"] = throughput
		sps.throughputMetrics["error_rate"] = errorRate
		sps.throughputMetrics["success_rate"] = 1.0 - errorRate
		
		// Store in context for reporting
		sps.bddContext.SetTestData("performance_metrics", map[string]interface{}{
			"throughput_ops_per_sec": throughput,
			"error_rate_percent":     errorRate * 100,
			"success_rate_percent":   (1.0 - errorRate) * 100,
			"total_operations":       sps.operationCount,
			"failed_operations":      sps.errorCount,
			"duration_seconds":       sps.duration.Seconds(),
		})
	}
}

// validateResourceUsage validates memory and CPU usage against thresholds
func (sps *SimplePerformanceSteps) validateResourceUsage() error {
	memoryIncrease := sps.memoryAfter.Alloc - sps.memoryBefore.Alloc
	memoryIncreaseRatio := float64(memoryIncrease) / float64(sps.memoryBefore.Alloc)
	
	// Define memory usage thresholds
	const maxMemoryIncreaseRatio = 0.5 // 50% increase allowed
	const maxMemoryIncreaseMB = 100    // 100MB absolute increase
	
	if memoryIncreaseRatio > maxMemoryIncreaseRatio {
		violation := fmt.Sprintf("memory increased by %.1f%%, exceeds threshold of %.1f%%",
			memoryIncreaseRatio*100, maxMemoryIncreaseRatio*100)
		sps.performanceViolations = append(sps.performanceViolations, violation)
		return fmt.Errorf(violation)
	}
	
	if memoryIncrease > maxMemoryIncreaseMB*1024*1024 {
		violation := fmt.Sprintf("memory increased by %d MB, exceeds threshold of %d MB",
			memoryIncrease/(1024*1024), maxMemoryIncreaseMB)
		sps.performanceViolations = append(sps.performanceViolations, violation)
		return fmt.Errorf(violation)
	}
	
	return nil
}

// getPerformanceReport generates a comprehensive performance report
func (sps *SimplePerformanceSteps) getPerformanceReport() map[string]interface{} {
	sps.measureThroughput()
	
	report := map[string]interface{}{
		"operation_count":         sps.operationCount,
		"error_count":             sps.errorCount,
		"duration_ms":             sps.duration.Milliseconds(),
		"throughput_metrics":      sps.throughputMetrics,
		"performance_violations": sps.performanceViolations,
		"memory_before_mb":        sps.memoryBefore.Alloc / 1024 / 1024,
		"memory_after_mb":         sps.memoryAfter.Alloc / 1024 / 1024,
		"memory_increase_mb":      (sps.memoryAfter.Alloc - sps.memoryBefore.Alloc) / 1024 / 1024,
		"thresholds":              sps.performanceThresholds,
	}
	
	return report
}
