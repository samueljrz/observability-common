#!/bin/bash

echo "🚀 OpenTelemetry Collector Setup"
echo "================================"

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running"
    echo ""
    echo "Please start Docker first:"
    echo "  - On macOS: Open Docker Desktop application"
    echo "  - On Linux: sudo systemctl start docker"
    echo "  - On Windows: Start Docker Desktop"
    echo ""
    echo "Or if using Colima:"
    echo "  colima start"
    echo ""
    exit 1
fi

echo "✅ Docker is running"

# Check if docker-compose is available
if ! command -v docker-compose > /dev/null 2>&1; then
    echo "❌ docker-compose is not installed"
    echo "Please install docker-compose first"
    exit 1
fi

echo "✅ docker-compose is available"

# Start the collector
echo "🔧 Starting OpenTelemetry collector..."
docker-compose up -d

if [ $? -eq 0 ]; then
    echo "✅ Collector started successfully!"
    echo ""
    echo "📊 Collector is now running on:"
    echo "  - OTLP gRPC: localhost:4317"
    echo "  - OTLP HTTP: localhost:4318"
    echo "  - Syslog: localhost:54526"
    echo ""
    echo "🔍 To view collector logs:"
    echo "  docker-compose logs -f otel-collector"
    echo ""
    echo "🚀 To run the example:"
    echo "  cd .. && go run example/collector/main.go"
    echo ""
    echo "🛑 To stop the collector:"
    echo "  docker-compose down"
else
    echo "❌ Failed to start collector"
    exit 1
fi 