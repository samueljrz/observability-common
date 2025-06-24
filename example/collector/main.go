package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	obs "github.com/garden/observability-commons"
	"github.com/garden/observability-commons/config"
)

const (
	collectorServiceName    = "collector-demo"
	collectorServiceVersion = "1.0.0"
	collectorSearchIndex    = "observability-collector-demo"
)

type paymentError struct {
	err error
}

func (e paymentError) Error() string {
	return e.err.Error()
}

func main() {
	fmt.Println("🚀 Starting Observability Collector Demo")
	fmt.Println("📡 This example sends data to the OpenTelemetry collector")
	fmt.Println("🔍 Check the collector logs to see the data flowing through")
	fmt.Println()

	observabilityClient, err := obs.NewObservability(config.Config{
		Service: config.Service{
			Name:    collectorServiceName,
			Version: collectorServiceVersion,
		},
		Mode:          config.Debug,
		SearchIndex:   collectorSearchIndex,
		FlushInterval: 2 * time.Second,
		Timeout:       10 * time.Second,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to create observability client: %v", err))
	}
	defer observabilityClient.Close()

	fmt.Println("✅ Observability client created successfully")
	fmt.Println("📊 Mode: Debug (sending to collector)")
	fmt.Println("🔄 Flush interval: 2 seconds")
	fmt.Println()

	// Example 1: Basic logging
	fmt.Println("📝 Example 1: Basic logging")
	for i := 0; i < 3; i++ {
		observabilityClient.Info(
			"order-service",
			"create-order",
			fmt.Sprintf("Order %d created successfully", i+1),
			map[string]string{
				"order_id":  fmt.Sprintf("order-%d", i+1),
				"amount":    fmt.Sprintf("%.2f", 99.99+float64(i)),
				"currency":  "USD",
				"customer":  fmt.Sprintf("customer-%d", i+1),
				"iteration": fmt.Sprintf("%d", i+1),
			},
		)
		time.Sleep(500 * time.Millisecond)
	}

	// Example 2: Error logging
	fmt.Println("❌ Example 2: Error logging")
	observabilityClient.Error(
		"payment-service",
		"process-payment",
		"Payment processing failed",
		paymentError{err: errors.New("insufficient funds")},
		map[string]string{
			"order_id":    "order-123",
			"payment_id":  "pay-456",
			"error_code":  "INSUFFICIENT_FUNDS",
			"retry_count": "3",
		},
	)

	// Example 3: Tracing
	fmt.Println("🔍 Example 3: Distributed tracing")
	ctx, span := observabilityClient.StartSpan(context.Background(), "process-order-workflow")
	defer span.End()

	// Add events to the span
	observabilityClient.AddEvent(ctx, "order-validated", map[string]string{
		"validation_time_ms": "150",
		"items_count":        "3",
	})

	time.Sleep(200 * time.Millisecond)

	observabilityClient.AddEvent(ctx, "payment-processed", map[string]string{
		"payment_method":     "credit_card",
		"processing_time_ms": "300",
	})

	time.Sleep(100 * time.Millisecond)

	observabilityClient.AddEvent(ctx, "inventory-updated", map[string]string{
		"items_reserved": "3",
		"warehouse":      "us-east-1",
	})

	// Example 4: Metrics
	fmt.Println("📊 Example 4: Metrics collection")
	ctx = context.Background()

	// Counter metrics
	observabilityClient.SystemMetricCounter(ctx, "orders.created", 1, map[string]string{
		"region": "us-east-1",
		"source": "web",
	})

	observabilityClient.SystemMetricCounter(ctx, "payments.processed", 1, map[string]string{
		"payment_method": "credit_card",
		"status":         "success",
	})

	// Histogram metrics
	processingTime := rand.Float64() * 1000
	observabilityClient.SystemMetricHistogram(ctx, "order.processing_time_ms", processingTime, map[string]string{
		"service": "order-service",
		"region":  "us-east-1",
	})

	// Gauge metrics
	activeOrders := rand.Int63n(100)
	observabilityClient.SystemMetricGauge(ctx, "orders.active", activeOrders, map[string]string{
		"region": "us-east-1",
	})

	fmt.Println("✅ All examples completed!")
	fmt.Println()
	fmt.Println("📋 What to check:")
	fmt.Println("   1. Look at the collector logs (docker logs otel-collector)")
	fmt.Println("   2. You should see structured logs, metrics, and traces")
	fmt.Println("   3. All data is being processed through the collector pipeline")
	fmt.Println()
	fmt.Println("🔄 Starting continuous data generation...")
	fmt.Println("   Press Enter to stop")

	// Start background data generation
	go generateContinuousData(observabilityClient)

	// Wait for user input to stop
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	fmt.Println("👋 Shutting down...")
}

func generateContinuousData(observabilityClient obs.Observability) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	counter := 0
	for range ticker.C {
		counter++

		// Generate random metrics
		ctx := context.Background()

		// Random processing time
		processingTime := rand.Float64() * 2000
		observabilityClient.SystemMetricHistogram(ctx, "api.request.duration_ms", processingTime, map[string]string{
			"endpoint": "/api/orders",
			"method":   "POST",
		})

		// Random active connections
		activeConnections := rand.Int63n(500)
		observabilityClient.SystemMetricGauge(ctx, "connections.active", activeConnections, map[string]string{
			"service": "api-gateway",
		})

		// Increment request counter
		observabilityClient.SystemMetricCounter(ctx, "api.requests.total", 1, map[string]string{
			"endpoint": "/api/orders",
			"status":   "200",
		})

		// Log some activity
		if counter%5 == 0 {
			observabilityClient.Info(
				"api-gateway",
				"handle-request",
				fmt.Sprintf("Processed %d requests", counter),
				map[string]string{
					"request_count":  fmt.Sprintf("%d", counter),
					"uptime_minutes": fmt.Sprintf("%d", counter/20),
				},
			)
		}

		// Occasionally log an error
		if counter%7 == 0 {
			observabilityClient.Warn(
				"api-gateway",
				"rate-limit",
				"Rate limit approaching",
				errors.New("high request volume"),
				map[string]string{
					"requests_per_minute": fmt.Sprintf("%d", rand.Intn(1000)),
					"limit":               "1000",
				},
			)
		}
	}
}
