# observability-commons

Modern observability library for projects written in Golang. It provides unified logging, metrics, and tracing capabilities using OpenTelemetry standards.

## Features

- **Unified Observability**: Single interface for logs, metrics, and tracing
- **OTLP-based Logging**: Modern OpenTelemetry logging instead of syslog
- **OpenTelemetry Integration**: Native support for metrics and tracing
- **Flexible Configuration**: Multiple deployment modes for different environments
- **Evidence-free Design**: Simplified logging without evidence complexity

## Quick Start

### Prerequisites
- Go 1.18 or later
- Git

### Installation
```bash
git clone <repository-url>
cd observability-commons
go mod download
```

## How to Run

### 1. **Using Make Commands (Recommended)**

The easiest way to run and test the library is using the provided Makefile:

```bash
# Test if the library can be built and run
make build-test

# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run the example application
make example

# Run benchmarks
make bench

# Format code
make fmt

# Run linter
make lint

# Install dependencies
make deps

# Clean build artifacts
make clean

# Verify everything works (build-test + test + example)
make verify
```

### 2. **Manual Go Commands**

You can also run the library directly with Go commands:

```bash
# Run all tests
go test -v ./...

# Run specific test
go test -v -run TestBuildAndRun

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run benchmarks
go test -benchmem -run=^$ -bench .

# Run the example
go run example/simple.go

# Build the library
go build -o bin/observability-commons .

# Format code
go fmt ./...

# Download dependencies
go mod download
```

### 3. **Running the Example Application**

The example demonstrates all features of the library:

```bash
# Run the example (will run indefinitely until you press Enter)
go run example/simple.go

# Or use make
make example
```

The example will:
- Send structured logs with component/operation information
- Demonstrate error logging with stack traces
- Show tracing with spans and events
- Send various types of metrics (histograms, gauges, counters)
- Display metrics output in JSON format

### 4. **Running Tests**

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific test file
go test -v observability_test.go

# Run specific test function
go test -v -run TestObservabilityLogs_NoError

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run benchmarks
go test -benchmem -run=^$ -bench .
```

### 5. **Running with Different Modes**

The library supports different logging modes. You can test them by modifying the example:

```bash
# Edit example/simple.go to change the mode
# Available modes: Noop, Local, Debug, Development, Production

# Noop mode (no output)
Mode: config.Noop

# Local mode (stdout)
Mode: config.Local

# Debug mode (OTLP collector)
Mode: config.Debug
```

### 6. **Running with Docker (Collector)**

To test with the OpenTelemetry collector:

```bash
# Start the collector
make start-collector

# Run the collector example (sends data to collector)
make collector-example

# View collector logs
make collector-logs

# Stop the collector
make stop-collector

# Or run the full demo (start collector + run example)
make collector-demo
```

**Manual collector setup:**
```bash
# Start the collector
cd collector
docker-compose up -d

# Run the example with Debug mode
go run example/collector/main.go

# View collector logs
docker-compose logs -f otel-collector
```

The collector example demonstrates:
- Sending structured logs to the collector
- Distributed tracing with spans and events
- Metrics collection (counters, histograms, gauges)
- Continuous data generation
- Real-time data flow through the collector pipeline

### 7. **Integration Testing**

```bash
# Run integration tests
go test -v -tags=integration ./...

# Run with specific environment variables
GARDEN_STACK=test go test -v ./...
```

### 8. **Performance Testing**

```bash
# Run benchmarks
go test -benchmem -run=^$ -bench .

# Run benchmarks with specific functions
go test -benchmem -run=^$ -bench BenchmarkObservabilityLogs

# Run benchmarks with profiling
go test -benchmem -run=^$ -bench . -cpuprofile=cpu.prof -memprofile=mem.prof
```

### 9. **Development Workflow**

```bash
# 1. Install dependencies
go mod download

# 2. Format code
go fmt ./...

# 3. Run linter
golangci-lint run

# 4. Run tests
go test -v ./...

# 5. Run example
go run example/simple.go

# 6. Build
go build -o bin/observability-commons .
```

### 10. **Continuous Integration**

The library includes a `.drone.yml` file for CI/CD:

```bash
# Run CI locally (if you have drone CLI)
drone exec

# Or run the same commands manually
make verify
```

## Expected Output

When running the example, you should see:

1. **Structured JSON logs** with component/operation information
2. **Metrics output** in JSON format showing histograms, gauges, and counters
3. **Console messages** indicating what's happening

Example log output:
```json
{
  "level": "info",
  "message": "Order processed successfully",
  "service.name": "weeb-app",
  "service.version": "0.41.7",
  "host.name": "your-hostname",
  "component": "order-service",
  "operation": "process-order",
  "timestamp": "2025-06-14T01:28:34.850-0300",
  "order_id": "order-0",
  "amount": "99.99",
  "currency": "USD"
}
```

## Troubleshooting

### Common Issues

1. **Import errors**: Make sure you're using the correct module name `github.com/garden/observability-commons`

2. **Missing dependencies**: Run `go mod download` to install all dependencies

3. **Build errors**: Ensure you're using Go 1.18 or later

4. **Test failures**: Check that all required environment variables are set

### Getting Help

- Check the example code in `example/simple.go`
- Review the test files for usage patterns
- Look at the configuration options in `config/config.go`

## Collector

The default export is to OTLP through OpenTelemetry collectors. If necessary, you can run a simple collector in your machine to simulate the data going through the OTLP receiver.

You can find a basic configuration inside the `/collector` folder, and so as to run it, you can run `docker-compose up`. You also need to set the logging mode to `Debug`, so the endpoint will be set to the OTLP port.

## Logging modes

This library has five modes of logging: `Noop`, `Debug`, `Local`, `Development`, and `Production`.

### _Noop_ mode

It can be used whenever you don't want to send the logs anywhere, not even to `stdout`.

### _Debug_ mode

It's used whenever you're running the collector locally, so it will send the logs to `localhost:54526`, accordingly to what is configured on the `/collector` folder.

### _Local_ mode

It's used whenever you want to send the logs to `stdout`.

### _Development_ mode

When you use this mode, the logs will be sent to the development collector. It's important to notice that in this library, whenever using this mode, the search index is automatically changed to append `_dev` as a suffix.

### _Production_ mode

By using this mode, the logs will be sent to the production collector.

## Usage

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
    observabilityClient, err := obs.NewObservability(config.Config{
        Service: config.Service{
            Name:    "my-app",
            Version: "1.0.0",
        },
        Mode:          config.Local,
        FlushInterval: 5 * time.Second,
        Timeout:       10 * time.Second,
    })
    if err != nil {
        panic(err)
    }
    defer observabilityClient.Close()

    // Logging
    observabilityClient.Info("order-service", "create-order", "Order created successfully", map[string]string{
        "amount": "99.99",
        "currency": "USD",
    })

    // Tracing
    ctx, span := observabilityClient.StartSpan(context.Background(), "process-order")
    defer span.End()
    
    observabilityClient.AddEvent(ctx, "payment-processed", map[string]string{
        "payment_method": "credit_card",
    })

    // Metrics
    observabilityClient.SystemMetricCounter(ctx, "orders.created", 1, map[string]string{
        "region": "us-east-1",
    })
}
```

## FAQ

- How to skip logging when executing unit tests?
  A: You can use the log mode `Noop`, which executes the logging pipeline, but it doesn't write anything neither to disk nor perform any IO.

- What's the purpose and best way to use the `Component` and `Operation` parameters?
  A: You can use at your will. Usually `Component` refers to the service or module name (eg. "order-service", "payment-service") and `Operation` refers to the specific action being performed (eg. "create-order", "process-payment").

- Not receiving data in the data source.
  A: First, use the library in local mode, to ensure that the logs are appearing in `stdout`. If it works on local mode, but you still cannot see the logs on your observability platform, check your collector configuration or ask for help.
