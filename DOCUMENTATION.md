# Observability Commons Library - Complete Documentation

## Table of Contents
1. [Project Overview](#project-overview)
2. [Folder Structure](#folder-structure)
3. [Core Components](#core-components)
4. [Package Documentation](#package-documentation)
5. [Configuration](#configuration)
6. [Usage Examples](#usage-examples)
7. [Architecture](#architecture)
8. [API Reference](#api-reference)

## Project Overview

The Observability Commons Library is a modern Go library that provides unified logging, metrics, and tracing capabilities using OpenTelemetry standards. It replaces traditional syslog-based logging with OTLP (OpenTelemetry Protocol) and provides a clean, component-based interface.

### Key Features
- **Unified Interface**: Single API for logs, metrics, and tracing
- **OTLP-based**: Modern OpenTelemetry protocol instead of syslog
- **Component-based**: Intuitive component/operation logging model
- **Multiple Modes**: Local, Debug, Development, Production, and Noop modes
- **OpenTelemetry Integration**: Native support for metrics and tracing
- **Evidence-free**: Simplified design without complex evidence handling

## Folder Structure

```
observability-commons/
├── 📁 Root Files
│   ├── README.md              # Main project documentation
│   ├── DOCUMENTATION.md       # This comprehensive guide
│   ├── Makefile              # Build and development commands
│   ├── go.mod                # Go module definition
│   ├── go.sum                # Dependency checksums
│   ├── .gitignore           # Git ignore patterns
│   └── observability.go     # Main observability interface
│
├── 📁 config/                # Configuration package
│   ├── config.go            # Configuration struct and validation
│   └── mode.go              # Logging mode definitions
│
├── 📁 log/                   # Logging package
│   ├── log.go               # Logger interface definition
│   ├── model.go             # Log entry data structures
│   └── otlp.go              # OTLP-based logger implementation
│
├── 📁 metrics/               # Metrics package
│   ├── metrics.go           # Metrics interface and implementation
│   └── client.go            # Metrics client utilities
│
├── 📁 trace/                 # Tracing package
│   └── trace.go             # OpenTelemetry tracing implementation
│
├── 📁 util/                  # Utility functions
│   ├── error.go             # Error handling utilities
│   ├── error_test.go        # Error utility tests
│   ├── fields.go            # Field processing utilities
│   └── hash.go              # Hash generation utilities
│
├── 📁 example/               # Example applications
│   ├── simple.go            # Basic usage example
│   └── collector/           # Collector integration example
│       └── main.go          # Collector demo application
│
└── 📁 collector/             # OpenTelemetry collector setup
    ├── README.md            # Collector-specific documentation
    ├── setup.sh             # Collector setup script
    ├── docker-compose.yml   # Docker Compose configuration
    └── otel-collector-config.yml  # Collector configuration
```

## Core Components

### 1. Main Interface (`observability.go`)

The main entry point that provides a unified interface for all observability features.

**Key Components:**
- `Observability` interface: Defines all logging, metrics, and tracing methods
- `ObservabilityClient` struct: Main implementation
- `NewObservability()` function: Factory function to create instances

**Methods:**
- **Logging**: `Debug()`, `Info()`, `Warn()`, `Error()`, `Fatal()`
- **Tracing**: `StartSpan()`, `AddEvent()`, `SetAttributes()`
- **Metrics**: `SystemMetricHistogram()`, `SystemMetricCounter()`, `SystemMetricGauge()`
- **Resource Management**: `Close()`

### 2. Configuration (`config/`)

Handles all configuration aspects of the library.

#### `config.go`
- `Config` struct: Main configuration container
- `Service` struct: Service name and version
- `Ensure()` method: Validates and sets default values
- `GetHostname()` and `GetSearchIndex()` methods: Utility getters

#### `mode.go`
- Defines logging modes: `Noop`, `Local`, `Debug`, `Development`, `Production`
- Each mode determines how data is processed and where it's sent

### 3. Logging (`log/`)

Modern OTLP-based logging implementation.

#### `log.go`
- `Logger` interface: Defines logging contract
- `Debug()`, `Info()`, `Warn()`, `Error()`, `Fatal()` methods

#### `model.go`
- `Entry` struct: Log entry data structure
- Contains: `Component`, `Operation`, `Message`, `Err`, `Fields`
- Internal fields for stacktrace handling

#### `otlp.go`
- `OTLPLogger` struct: OTLP-based logger implementation
- Uses Zap logger for structured JSON output
- Supports different modes (Noop, Local, Debug, etc.)
- Generates OTLP-compatible log fields

### 4. Metrics (`metrics/`)

OpenTelemetry metrics implementation.

#### `metrics.go`
- `Meter` interface: Defines metrics contract
- `OtelMeter` struct: OpenTelemetry implementation
- Supports: Histograms, Gauges, Counters
- Automatic resource attributes and default fields

#### `client.go`
- Metrics client utilities and helpers
- Endpoint configuration for different modes

### 5. Tracing (`trace/`)

OpenTelemetry distributed tracing implementation.

#### `trace.go`
- `Tracer` interface: Defines tracing contract
- `OtelTracer` struct: OpenTelemetry implementation
- `Span` interface: Span operations
- Supports: Span creation, events, attributes
- No-op implementation for compatibility

### 6. Utilities (`util/`)

Helper functions and utilities.

#### `error.go`
- Error handling utilities
- Error wrapping and processing

#### `fields.go`
- Field processing utilities
- `ExtraFields` type for structured data

#### `hash.go`
- `MD5Hash()` function: Generates MD5 hashes for stacktraces

#### `error_test.go`
- Unit tests for error utilities

## Package Documentation

### Configuration Package (`config/`)

**Purpose**: Manages all configuration aspects of the observability library.

**Key Features:**
- Service identification (name, version)
- Logging mode selection
- Endpoint configuration
- Default field management
- Hostname and search index handling

**Usage:**
```go
cfg := config.Config{
    Service: config.Service{
        Name:    "my-service",
        Version: "1.0.0",
    },
    Mode:          config.Local,
    SearchIndex:   "my-service-logs",
    FlushInterval: 5 * time.Second,
    Timeout:       10 * time.Second,
}
```

### Logging Package (`log/`)

**Purpose**: Provides structured, OTLP-compatible logging.

**Key Features:**
- Component/operation-based logging
- Structured JSON output
- Stacktrace handling
- OTLP-compatible field generation
- Multiple output modes

**Usage:**
```go
logger, err := log.NewOTLPLogger(cfg)
logger.Info(&log.Entry{
    Component:  "order-service",
    Operation:  "create-order",
    Message:    "Order created successfully",
    Fields:     map[string]string{"order_id": "123"},
})
```

### Metrics Package (`metrics/`)

**Purpose**: Provides OpenTelemetry metrics collection.

**Key Features:**
- Counter, histogram, and gauge metrics
- Automatic resource attributes
- Default field injection
- Multiple export modes
- Batch processing

**Usage:**
```go
meter, err := metrics.NewOtelMeter(cfg)
meter.DefaultCounter(ctx, "orders.created", 1, fields)
meter.DefaultHistogram(ctx, "order.processing_time", 150.5, fields)
meter.DefaultGauge(ctx, "active.orders", 42, fields)
```

### Tracing Package (`trace/`)

**Purpose**: Provides distributed tracing capabilities.

**Key Features:**
- Span creation and management
- Event and attribute handling
- OpenTelemetry integration
- No-op fallback for compatibility

**Usage:**
```go
tracer, err := trace.NewTracer(cfg)
ctx, span := tracer.StartSpan(context.Background(), "process-order")
defer span.End()
tracer.AddEvent(ctx, "payment-processed", attributes)
```

### Utilities Package (`util/`)

**Purpose**: Provides helper functions and utilities.

**Key Features:**
- Error handling utilities
- Field processing
- Hash generation
- Type conversions

**Usage:**
```go
hash := util.MD5Hash([]byte("stacktrace"))
fields := util.ExtraFields{"key": "value"}
```

## Configuration

### Logging Modes

1. **Noop**: No output, useful for testing
2. **Local**: Output to stdout (JSON format)
3. **Debug**: Send to local collector (localhost:4317)
4. **Development**: Send to development collector with `_dev` suffix
5. **Production**: Send to production collector

### Configuration Options

```go
type Config struct {
    Service       Service              // Service identification
    Mode          Mode                 // Logging mode
    SearchIndex   string               // Search index name
    FlushInterval time.Duration        // Metrics flush interval
    Timeout       time.Duration        // Request timeout
    Port          string               // Collector port (default: 80)
    DefaultFields *map[string]string   // Default fields for all data
}
```

## Usage Examples

### Basic Usage

```go
package main

import (
    "context"
    "time"
    
    obs "github.com/garden/observability-commons"
    "github.com/garden/observability-commons/config"
)

func main() {
    // Create observability client
    client, err := obs.NewObservability(config.Config{
        Service: config.Service{
            Name:    "my-service",
            Version: "1.0.0",
        },
        Mode:          config.Local,
        FlushInterval: 5 * time.Second,
        Timeout:       10 * time.Second,
    })
    if err != nil {
        panic(err)
    }
    defer client.Close()

    // Logging
    client.Info("order-service", "create-order", "Order created", map[string]string{
        "order_id": "123",
        "amount":   "99.99",
    })

    // Tracing
    ctx, span := client.StartSpan(context.Background(), "process-order")
    defer span.End()
    
    client.AddEvent(ctx, "payment-processed", map[string]string{
        "payment_method": "credit_card",
    })

    // Metrics
    client.SystemMetricCounter(ctx, "orders.created", 1, map[string]string{
        "region": "us-east-1",
    })
}
```

### Advanced Usage with Collector

```go
// Use Debug mode to send to collector
client, err := obs.NewObservability(config.Config{
    Service: config.Service{
        Name:    "collector-demo",
        Version: "1.0.0",
    },
    Mode:          config.Debug,  // Sends to localhost:4317
    SearchIndex:   "demo-logs",
    FlushInterval: 2 * time.Second,
    Timeout:       10 * time.Second,
})
```

## Architecture

### Component Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Application   │    │  Observability  │    │   Collector     │
│                 │    │     Client      │    │                 │
│ ┌─────────────┐ │    │ ┌─────────────┐ │    │ ┌─────────────┐ │
│ │   Logging   │◄┼────┼►│   Logger    │◄┼────┼►│   OTLP      │ │
│ └─────────────┘ │    │ └─────────────┘ │    │ └─────────────┘ │
│ ┌─────────────┐ │    │ ┌─────────────┐ │    │ ┌─────────────┐ │
│ │   Metrics   │◄┼────┼►│    Meter    │◄┼────┼►│   Metrics   │ │
│ └─────────────┘ │    │ └─────────────┘ │    │ └─────────────┘ │
│ ┌─────────────┐ │    │ ┌─────────────┐ │    │ ┌─────────────┐ │
│ │   Tracing   │◄┼────┼►│   Tracer    │◄┼────┼►│   Traces    │ │
│ └─────────────┘ │    │ └─────────────┘ │    │ └─────────────┘ │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Data Flow

1. **Application** calls observability methods
2. **ObservabilityClient** processes and formats data
3. **Individual components** (Logger, Meter, Tracer) handle specific data types
4. **OTLP exporters** send data to collector or local output
5. **Collector** processes and forwards data to final destinations

### Mode Processing

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Noop      │    │   Local     │    │   Debug     │
│             │    │             │    │             │
│ No output   │    │ stdout      │    │ localhost   │
│             │    │ JSON        │    │ :4317       │
└─────────────┘    └─────────────┘    └─────────────┘
       │                   │                   │
       └───────────────────┼───────────────────┘
                           │
                    ┌─────────────┐
                    │Development/ │
                    │Production   │
                    │             │
                    │ Remote      │
                    │ collectors  │
                    └─────────────┘
```

## API Reference

### Observability Interface

```go
type Observability interface {
    // Logging methods
    Debug(component, operation, message string, fields map[string]string)
    Info(component, operation, message string, fields map[string]string)
    Warn(component, operation, message string, err error, fields map[string]string)
    Error(component, operation, message string, err error, fields map[string]string)
    Fatal(component, operation, message string, err error, fields map[string]string)

    // Tracing methods
    StartSpan(ctx context.Context, name string, opts ...trace.SpanOption) (context.Context, trace.Span)
    AddEvent(ctx context.Context, name string, attributes map[string]string)
    SetAttributes(ctx context.Context, attributes map[string]string)

    // Metrics methods
    SystemMetricHistogram(ctx context.Context, metricName string, value float64, fields map[string]string) error
    SystemMetricCounter(ctx context.Context, metricName string, value int64, fields map[string]string) error
    SystemMetricGauge(ctx context.Context, metricName string, value int64, fields map[string]string) error

    // Resource management
    Close() error
}
```

### Configuration Struct

```go
type Config struct {
    Service       Service              // Service identification
    Mode          Mode                 // Logging mode
    SearchIndex   string               // Search index name
    FlushInterval time.Duration        // Metrics flush interval
    Timeout       time.Duration        // Request timeout
    Port          string               // Collector port
    DefaultFields *map[string]string   // Default fields
}
```

### Log Entry Struct

```go
type Entry struct {
    Component  string            // Service component
    Operation  string            // Operation being performed
    Message    string            // Log message
    Err        error             // Error (if any)
    Fields     map[string]string // Additional fields
}
```

## Development

### Building

```bash
# Build the library
make build

# Run tests
make test

# Run with coverage
make test-coverage

# Format code
make fmt

# Run linter
make lint
```

### Testing

```bash
# Run all tests
go test ./...

# Run specific test
go test -v -run TestObservabilityLogs_NoError

# Run benchmarks
go test -benchmem -run=^$ -bench .
```

### Examples

```bash
# Run basic example
make example

# Run collector example
make collector-demo

# Or manually
go run example/simple.go
go run example/collector/main.go
```

## Best Practices

### 1. Component and Operation Naming

- **Component**: Use service or module names (e.g., "order-service", "payment-service")
- **Operation**: Use specific actions (e.g., "create-order", "process-payment")

### 2. Field Usage

- Use consistent field names across your application
- Include relevant context (IDs, timestamps, user info)
- Avoid sensitive information in fields

### 3. Error Handling

- Always include errors in Error/Warn/Fatal calls
- Use structured error types when possible
- Include relevant context in error fields

### 4. Performance

- Use appropriate flush intervals for your use case
- Consider using Noop mode for tests
- Batch operations when possible

### 5. Configuration

- Use environment-specific modes
- Set appropriate timeouts and intervals
- Include default fields for common context

## Troubleshooting

### Common Issues

1. **Import errors**: Ensure correct module name `github.com/garden/observability-commons`
2. **Docker issues**: Check Docker is running for collector examples
3. **Port conflicts**: Verify ports 4317, 4318, 54526 are available
4. **Configuration errors**: Validate config with `Ensure()` method

### Debugging

1. **Check logs**: Use `make collector-logs` for collector debugging
2. **Test connectivity**: Use `telnet localhost 4317` for OTLP testing
3. **Verify configuration**: Check config validation in logs
4. **Mode testing**: Try different modes to isolate issues

This documentation provides a complete overview of the observability commons library, its structure, and usage patterns. For specific implementation details, refer to the individual package documentation and source code. 