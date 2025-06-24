.PHONY: build test clean build-test

# Build the observability commons library
build:
	go build -o bin/observability-commons .

# Test if the library can be built and run
build-test:
	@echo "Testing if the library can be built and run..."
	@go test -v -run TestBuildAndRun
	@echo "✅ Build test passed!"

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Run the example
example:
	go run example/simple.go

# Run the collector example
collector-example:
	@echo "🚀 Starting collector example..."
	@echo "📡 Make sure the collector is running: make start-collector"
	@echo ""
	go run example/collector/main.go

# Start the OpenTelemetry collector
start-collector:
	@echo "🔧 Starting OpenTelemetry collector..."
	@cd collector && ./setup.sh

# Stop the OpenTelemetry collector
stop-collector:
	@echo "🛑 Stopping OpenTelemetry collector..."
	@cd collector && docker-compose down
	@echo "✅ Collector stopped"

# View collector logs
collector-logs:
	@echo "📋 Collector logs (following):"
	@cd collector && docker-compose logs -f

# View collector logs (last 50 lines)
collector-logs-tail:
	@echo "📋 Last 50 lines of collector logs:"
	@cd collector && docker-compose logs --tail 50

# View collector logs with timestamps
collector-logs-timestamps:
	@echo "📋 Collector logs with timestamps:"
	@cd collector && docker-compose logs -t

# View collector logs since specific time (last 10 minutes)
collector-logs-recent:
	@echo "📋 Collector logs from last 10 minutes:"
	@cd collector && docker-compose logs --since 10m

# View collector logs for specific container by ID
collector-logs-id:
	@echo "📋 Enter collector container ID:"
	@read container_id; docker logs -f $$container_id

# View collector logs for specific container by ID (last 50 lines)
collector-logs-id-tail:
	@echo "📋 Enter collector container ID:"
	@read container_id; docker logs --tail 50 $$container_id

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Install dependencies
deps:
	go mod download

run_example:
	@go run ./example/simple.go

bench:
	@go test -benchmem -run=^$ -bench .

test_all: test bench

# Quick verification that everything works
verify: build-test test example
	@echo "✅ All verification tests passed!"

# Full collector demo
collector-demo: start-collector collector-example
	@echo "🎉 Collector demo completed!"