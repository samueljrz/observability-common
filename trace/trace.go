package trace

import (
	"context"
	"fmt"

	"github.com/garden/observability-commons/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const (
	instrumentationName = "github.com/garden/observability-commons"
)

type Span interface {
	End()
	AddEvent(name string, attributes map[string]string)
	SetAttributes(attributes map[string]string)
	SpanContext() trace.SpanContext
}

type SpanOption func()

type Tracer interface {
	StartSpan(ctx context.Context, name string, opts ...SpanOption) (context.Context, Span)
	AddEvent(ctx context.Context, name string, attributes map[string]string)
	SetAttributes(ctx context.Context, attributes map[string]string)
	Close() error
}

type OtelTracer struct {
	tracer trace.Tracer
	tp     *sdktrace.TracerProvider
	cfg    config.Config
}

func NewTracer(cfg config.Config) (*OtelTracer, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			attribute.String("service.name", cfg.Service.Name),
			attribute.String("service.version", cfg.Service.Version),
			attribute.String("host.name", cfg.GetHostname()),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	exporter := &noopExporter{}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	return &OtelTracer{
		tracer: tp.Tracer(instrumentationName),
		tp:     tp,
		cfg:    cfg,
	}, nil
}

func (t *OtelTracer) StartSpan(ctx context.Context, name string, opts ...SpanOption) (context.Context, Span) {
	spanCtx, span := t.tracer.Start(ctx, name)
	return spanCtx, &otelSpan{span: span}
}

func (t *OtelTracer) AddEvent(ctx context.Context, name string, attributes map[string]string) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
	}
}

func (t *OtelTracer) SetAttributes(ctx context.Context, attributes map[string]string) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
	}
}

func (t *OtelTracer) Close() error {
	if t.tp != nil {
		return t.tp.Shutdown(context.Background())
	}
	return nil
}

type otelSpan struct {
	span trace.Span
}

func (s *otelSpan) End() {
	s.span.End()
}

func (s *otelSpan) AddEvent(name string, attributes map[string]string) {
}

func (s *otelSpan) SetAttributes(attributes map[string]string) {
}

func (s *otelSpan) SpanContext() trace.SpanContext {
	return s.span.SpanContext()
}

type noopExporter struct{}

func (e *noopExporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	return nil
}

func (e *noopExporter) Shutdown(ctx context.Context) error {
	return nil
}
