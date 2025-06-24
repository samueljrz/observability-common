# OpenTelemetry Collector Setup

This directory contains the OpenTelemetry collector configuration for testing the observability commons library.

## Quick Start

### 1. **Automatic Setup (Recommended)**

```bash
# From the project root
make start-collector
make collector-example
```

### 2. **Manual Setup**

```bash
# Navigate to collector directory
cd collector

# Start the collector
./setup.sh

# Or manually with docker-compose
docker-compose up -d
```

## What's Included

### **Collector Configuration**
- **OTLP gRPC receiver**: `localhost:4317`
- **OTLP HTTP receiver**: `localhost:4318`
- **Syslog receiver**: `localhost:54526`
- **Logging exporter**: Outputs all data to stdout
- **Batch processor**: Batches data for efficiency

### **Docker Compose**
- OpenTelemetry Collector Contrib image
- Proper port mapping
- Volume mounting for configuration

## Usage

### **Start the Collector**
```bash
# From project root
make start-collector

# Or manually
cd collector
docker-compose up -d
```

### **Run the Example**
```bash
# From project root
make collector-example

# Or manually
go run example/collector/main.go
```

### **View Collector Logs**
```bash
# From project root
make collector-logs

# Or manually
cd collector
docker-compose logs -f otel-collector
```

### **Stop the Collector**
```bash
# From project root
make stop-collector

# Or manually
cd collector
docker-compose down
```

## What You'll See

When running the collector example, you should see:

### **In the Application Console:**
```
🚀 Starting Observability Collector Demo
📡 This example sends data to the OpenTelemetry collector
🔍 Check the collector logs to see the data flowing through

✅ Observability client created successfully
📊 Mode: Debug (sending to collector)
🔄 Flush interval: 2 seconds

📝 Example 1: Basic logging
❌ Example 2: Error logging
🔍 Example 3: Distributed tracing
📊 Example 4: Metrics collection
✅ All examples completed!
```

### **In the Collector Logs:**
```
ResourceLog #0
Resource labels:
     -> service.name: STRING(collector-demo)
     -> service.version: STRING(1.0.0)
     -> host.name: STRING(your-hostname)
InstrumentationLibraryLogs #0
InstrumentationLibrary
LogRecord #0
Timestamp: 2025-06-14 01:28:34.850 -0300
SeverityText: INFO
SeverityNumber: 9
Body: STRING(Order 1 created successfully)
Attributes:
     -> component: STRING(order-service)
     -> operation: STRING(create-order)
     -> order_id: STRING(order-1)
     -> amount: STRING(99.99)
     -> currency: STRING(USD)
```

## Configuration Details

### **Receivers**
- **OTLP gRPC**: Receives OpenTelemetry data via gRPC
- **OTLP HTTP**: Receives OpenTelemetry data via HTTP
- **Syslog**: Receives syslog messages (legacy support)

### **Processors**
- **Batch**: Batches data to improve performance

### **Exporters**
- **Logging**: Outputs all data to stdout for debugging
- **OTLP**: Forwards data to production collector (optional)

### **Service Pipelines**
- **Traces**: Processes distributed tracing data
- **Metrics**: Processes metrics data
- **Logs**: Processes log data

## Troubleshooting

### **Docker Not Running**
```bash
# On macOS
open -a Docker

# On Linux
sudo systemctl start docker

# With Colima
colima start
```

### **Port Already in Use**
```bash
# Check what's using the ports
lsof -i :4317
lsof -i :4318
lsof -i :54526

# Stop the collector first
docker-compose down
```

### **Permission Issues**
```bash
# Make setup script executable
chmod +x setup.sh
```

### **Collector Not Receiving Data**
1. Check if collector is running: `docker-compose ps`
2. Check collector logs: `docker-compose logs otel-collector`
3. Verify ports are exposed: `docker-compose port otel-collector 4317`
4. Test connectivity: `telnet localhost 4317`

## Advanced Configuration

### **Custom Endpoints**
Edit `otel-collector-config.yml` to change:
- Receiver endpoints
- Exporter destinations
- Processing pipelines

### **Environment Variables**
```bash
# Set custom environment
export OTEL_SERVICE_NAME=my-service
export OTEL_SERVICE_VERSION=1.0.0
```

### **Production Setup**
For production use:
1. Update exporter endpoints
2. Add authentication
3. Configure sampling
4. Set up monitoring for the collector itself

## Integration with Observability Commons

The collector example demonstrates:
- **Structured Logging**: Component/operation-based logs
- **Distributed Tracing**: Spans and events
- **Metrics Collection**: Counters, histograms, gauges
- **Real-time Processing**: Data flowing through collector pipeline
- **Batch Processing**: Efficient data handling
- **Multiple Protocols**: OTLP gRPC and HTTP support 