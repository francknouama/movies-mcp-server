# Reusable Packages (`pkg/`)

This directory contains reusable Go packages that can be used by external applications or shared across different projects. These packages follow Go best practices and have minimal dependencies on application-specific code.

## Available Packages

### 📊 `validation/` - Input Validation
Comprehensive validation framework with support for:
- **Basic validation rules**: required, min/max length, email, URL, regex patterns
- **Type validation**: alpha, numeric, alphanumeric, UUID, JSON
- **Custom validation**: MCP protocol validation, movie data validation
- **Struct tag validation**: Validate structs using field tags
- **Performance optimized**: Efficient validation pipeline

```go
import "movies-mcp-server/pkg/validation"

validator := validation.NewValidator()
validator.AddRule("email", validation.Required())
validator.AddRule("email", validation.Email())

err := validator.Validate(map[string]interface{}{
    "email": "user@example.com",
})
```

### 🚨 `errors/` - Structured Error Handling
Production-ready error handling with:
- **Standardized error codes** and severity levels
- **Structured error details** with context information
- **JSON-RPC compatibility** for API responses
- **Stack trace capture** for debugging
- **Error recovery utilities** with panic handling

```go
import "movies-mcp-server/pkg/errors"

err := errors.NewValidationError("Invalid input", "email", "user@invalid")
appErr := err.WithSeverity(errors.SeverityHigh).WithComponent("auth")
```

### ⏱️ `timeout/` - Timeout Management
Comprehensive timeout handling system:
- **Multi-layer timeouts**: request, database, image processing, health checks
- **Context-aware handling** with cancellation support
- **Circuit breaker pattern** for fault tolerance
- **Graceful shutdown** management
- **Timeout middleware** for automatic protection

```go
import "movies-mcp-server/pkg/timeout"

manager := timeout.NewManager(timeout.DefaultTimeoutConfig(), logger)
ctx, cancel := manager.WithRequestTimeout(context.Background())
defer cancel()
```

### 📝 `logging/` - Structured Logging
JSON-structured logging with:
- **Multiple log levels**: debug, info, warn, error
- **Context-aware logging** with request IDs
- **Specialized loggers** for different components
- **Performance metrics logging**
- **RFC3339 timestamps** for consistency

```go
import "movies-mcp-server/pkg/logging"

logger := logging.New(logging.LevelInfo)
logger.Info("user_action", "user_id", 123, "action", "login")
```

### 📈 `metrics/` - Performance Metrics
Real-time metrics collection:
- **Multiple metric types**: counters, gauges, histograms, timers
- **Built-in system metrics**: memory, goroutines
- **Request tracking** with automatic middleware
- **Configurable reporting** intervals
- **Thread-safe operations**

```go
import "movies-mcp-server/pkg/metrics"

metrics := metrics.NewMetrics(logger, 30*time.Second)
counter := metrics.NewCounter("requests_total", "Total requests")
metrics.IncCounter(counter)
```

### 🏥 `health/` - Health Checks
Comprehensive health monitoring:
- **Multiple health checkers**: database, memory, custom components
- **Readiness and liveness probes** for container orchestration
- **Detailed status reporting** with timing information
- **Configurable timeouts** for health checks
- **Extensible checker interface**

```go
import "movies-mcp-server/pkg/health"

manager := health.NewManager(logger, "1.0.0")
manager.RegisterChecker("database", health.NewDatabaseChecker(db))
healthStatus := manager.CheckAll(context.Background())
```

### 🖼️ `image/` - Image Processing
Image handling utilities:
- **Format validation**: JPEG, PNG, WebP support
- **Size validation** and constraints
- **Base64 encoding/decoding** for transport
- **MIME type detection** and validation
- **Thumbnail generation** (configurable)

```go
import "movies-mcp-server/pkg/image"

processor := image.NewImageProcessor(config)
err := processor.ValidateImage(imageData, "image/jpeg")
encoded := processor.EncodeToBase64(imageData, "image/jpeg")
```

## Design Principles

### 🔌 **Minimal Dependencies**
- Packages use only standard library and essential external dependencies
- No circular dependencies between pkg packages
- Clear separation between reusable utilities and application logic

### 📋 **Interface-Driven Design**
- Well-defined interfaces for extensibility
- Easy to mock and test
- Support for dependency injection

### 🧪 **Comprehensive Testing**
- Unit tests with >90% coverage
- Performance benchmarks for critical paths
- Integration examples and documentation

### 🔒 **Production Ready**
- Error handling and edge case coverage
- Performance optimizations
- Memory-efficient implementations
- Thread-safe operations where applicable

## Usage Guidelines

### ✅ **When to Use pkg/**
- Utilities that can be shared across multiple applications
- Generic functionality not tied to specific business logic
- Libraries that external consumers might find useful
- Code that follows Go standard library patterns

### ❌ **When to Use internal/**
- Application-specific business logic
- Models and types specific to your domain
- Configuration and setup code
- Database schemas and migrations

## Example: Using Multiple Packages Together

```go
package main

import (
    "context"
    "movies-mcp-server/pkg/logging"
    "movies-mcp-server/pkg/metrics"
    "movies-mcp-server/pkg/timeout"
    "movies-mcp-server/pkg/validation"
    "movies-mcp-server/pkg/health"
)

func main() {
    // Initialize logging
    logger := logging.New(logging.LevelInfo)
    
    // Setup metrics
    metrics := metrics.NewMetrics(logger, 30*time.Second)
    
    // Configure timeouts
    timeoutManager := timeout.NewManager(timeout.DefaultTimeoutConfig(), logger)
    
    // Setup validation
    validator := validation.NewRequestValidator()
    
    // Initialize health checks
    healthManager := health.NewManager(logger, "1.0.0")
    
    // Use together in your application
    ctx, cancel := timeoutManager.WithRequestTimeout(context.Background())
    defer cancel()
    
    timer := metrics.StartRequestTimer()
    defer metrics.FinishRequestTimer(timer)
    
    logger.Info("application_started", "version", "1.0.0")
}
```

This structure follows Go community best practices and makes it easy for other developers to consume and extend the functionality.