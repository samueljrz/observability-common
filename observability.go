package observability

import (
	"context"

	"github.com/garden/observability-commons/config"
	"github.com/garden/observability-commons/log"
	"github.com/garden/observability-commons/metrics"
	"github.com/garden/observability-commons/trace"
)

// Observability provides a unified interface for logging, metrics, and tracing
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

// ObservabilityClient is the main implementation of the Observability interface
type ObservabilityClient struct {
	logger log.Logger
	tracer trace.Tracer
	meter  metrics.Meter
}

// NewObservability creates a new observability client with OTLP-based logging and improved instrumentation
func NewObservability(cfg config.Config) (*ObservabilityClient, error) {
	err := cfg.Ensure()
	if err != nil {
		return nil, err
	}

	// Initialize OTLP-based logger instead of syslog
	logger, err := log.NewOTLPLogger(cfg)
	if err != nil {
		return nil, err
	}

	// Initialize tracer
	tracer, err := trace.NewTracer(cfg)
	if err != nil {
		return nil, err
	}

	// Initialize metrics
	meter, err := metrics.NewOtelMeter(cfg)
	if err != nil {
		return nil, err
	}

	return &ObservabilityClient{
		logger: logger,
		tracer: tracer,
		meter:  meter,
	}, nil
}

// Logging methods
func (obs *ObservabilityClient) Debug(component, operation, message string, fields map[string]string) {
	obs.logger.Debug(&log.Entry{
		Component: component,
		Operation: operation,
		Message:   message,
		Err:       nil,
		Fields:    fields,
	})
}

func (obs *ObservabilityClient) Info(component, operation, message string, fields map[string]string) {
	obs.logger.Info(&log.Entry{
		Component: component,
		Operation: operation,
		Message:   message,
		Err:       nil,
		Fields:    fields,
	})
}

func (obs *ObservabilityClient) Warn(component, operation, message string, err error, fields map[string]string) {
	obs.logger.Warn(&log.Entry{
		Component: component,
		Operation: operation,
		Message:   message,
		Err:       err,
		Fields:    fields,
	})
}

func (obs *ObservabilityClient) Error(component, operation, message string, err error, fields map[string]string) {
	obs.logger.Error(&log.Entry{
		Component: component,
		Operation: operation,
		Message:   message,
		Err:       err,
		Fields:    fields,
	})
}

func (obs *ObservabilityClient) Fatal(component, operation, message string, err error, fields map[string]string) {
	obs.logger.Fatal(&log.Entry{
		Component: component,
		Operation: operation,
		Message:   message,
		Err:       err,
		Fields:    fields,
	})
}

// Tracing methods
func (obs *ObservabilityClient) StartSpan(ctx context.Context, name string, opts ...trace.SpanOption) (context.Context, trace.Span) {
	return obs.tracer.StartSpan(ctx, name, opts...)
}

func (obs *ObservabilityClient) AddEvent(ctx context.Context, name string, attributes map[string]string) {
	obs.tracer.AddEvent(ctx, name, attributes)
}

func (obs *ObservabilityClient) SetAttributes(ctx context.Context, attributes map[string]string) {
	obs.tracer.SetAttributes(ctx, attributes)
}

// Metrics methods
func (obs *ObservabilityClient) SystemMetricHistogram(ctx context.Context, metricName string, value float64, fields map[string]string) error {
	return obs.meter.DefaultHistogram(ctx, metricName, value, fields)
}

func (obs *ObservabilityClient) SystemMetricCounter(ctx context.Context, metricName string, value int64, fields map[string]string) error {
	return obs.meter.DefaultCounter(ctx, metricName, value, fields)
}

func (obs *ObservabilityClient) SystemMetricGauge(ctx context.Context, metricName string, value int64, fields map[string]string) error {
	return obs.meter.DefaultGauge(ctx, metricName, value, fields)
}

// Close gracefully shuts down all observability components
func (obs *ObservabilityClient) Close() error {
	// Close logger
	if err := obs.logger.Close(); err != nil {
		return err
	}

	// Close tracer
	if err := obs.tracer.Close(); err != nil {
		return err
	}

	return nil
}
