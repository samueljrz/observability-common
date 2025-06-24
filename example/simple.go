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
	serviceName    = "web-app"
	serviceVersion = "0.41.7"
	searchIndex    = "observability-commons-example"
	metricPrefix   = "observability-commons-example"
)

type myCustomError struct {
	err error
}

func (customErr myCustomError) Error() string {
	return customErr.err.Error()
}

func main() {
	observabilityClient, err := obs.NewObservability(config.Config{
		Service: config.Service{
			Name:    serviceName,
			Version: serviceVersion,
		},
		Mode:          config.Local,
		SearchIndex:   searchIndex,
		FlushInterval: 5 * time.Second,
		Timeout:       10 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	defer observabilityClient.Close()

	// Example logging
	fmt.Println("Sending logs...")
	for i := 0; i < 5; i++ {
		observabilityClient.Info(
			"order-service",
			"process-order",
			"Order processed successfully",
			map[string]string{
				"order_id":  fmt.Sprintf("order-%d", i),
				"amount":    "99.99",
				"currency":  "USD",
				"iteration": fmt.Sprintf("%d", i),
			},
		)

		// Example error logging
		if i%2 == 0 {
			observabilityClient.Error(
				"order-service",
				"process-order",
				"Failed to process order",
				myCustomError{err: errors.New("payment failed")},
				map[string]string{
					"order_id": fmt.Sprintf("order-%d", i),
					"reason":   "payment_declined",
				},
			)
		}
	}

	// Example tracing
	fmt.Println("Demonstrating tracing...")
	ctx, span := observabilityClient.StartSpan(context.Background(), "process-orders")
	defer span.End()

	observabilityClient.AddEvent(ctx, "orders-processed", map[string]string{
		"count":  "5",
		"status": "completed",
	})

	// Example metrics
	fmt.Println("Sending metrics...")
	observabilityClient.SystemMetricCounter(ctx, "orders.processed", 5, map[string]string{
		"region": "us-east-1",
	})

	// Start background metrics collection
	go sendSystemMetrics(observabilityClient)

	fmt.Println("Observability example running. Press Enter to exit...")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
}

func sendSystemMetrics(observabilityClient obs.Observability) {
	for {
		delay := randomDelay(1, 1000)

		ctx := context.Background()

		// Send histogram metric
		fmt.Printf("Sending histogram metric: %s = %f\n", "processing_delay_ms", delay)
		observabilityClient.SystemMetricHistogram(ctx, "processing_delay_ms", delay, map[string]string{
			"service": "checkout",
		})

		// Send gauge metric
		fmt.Printf("Sending gauge metric: %s = %v\n", "active_orders", int64(delay/100))
		observabilityClient.SystemMetricGauge(ctx, "active_orders", int64(delay/100), map[string]string{
			"service": "checkout",
		})

		// Send counter metric
		fmt.Printf("Sending counter metric: %s = %v\n", "orders_processed_total", 1)
		observabilityClient.SystemMetricCounter(ctx, "orders_processed_total", 1, map[string]string{
			"service": "checkout",
		})

		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
}

func randomDelay(min, max int) float64 {
	rand.Seed(time.Now().UnixNano())
	return float64(rand.Intn(max-min+1) + min)
}
