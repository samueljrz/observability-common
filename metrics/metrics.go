package metrics

import (
	"context"
	"fmt"
	"os"

	"github.com/garden/observability-commons/config"
	"github.com/garden/observability-commons/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
)

const (
	instrumentationName = "github.com/garden/observability-commons"
)

var (
	defaultHistogramBuckets []float64
)

type Meter interface {
	DefaultHistogram(ctx context.Context, metricName string, value float64, fields util.ExtraFields) error
	DefaultGauge(ctx context.Context, metricName string, value int64, fields util.ExtraFields) error
	DefaultCounter(ctx context.Context, metricName string, value int64, fields util.ExtraFields) error
}

type OtelMeter struct {
	meter metric.Meter
	cfg   config.Config
}

func NewOtelMeter(cfg config.Config) (*OtelMeter, error) {
	ctx := context.Background()

	var exporter export.Exporter
	var err error
	switch cfg.Mode {
	case config.Noop:
		return &OtelMeter{
			meter: metric.NewNoopMeter(),
		}, nil
	case config.Local:
		exporter, err = stdoutmetric.New(stdoutmetric.WithPrettyPrint())
		if err != nil {
			return nil, fmt.Errorf("error creating otel exporter: %w", err)
		}
	case config.Debug, config.Development, config.Production:
		exporter, err = otlpmetric.New(ctx, newClient(cfg))
		if err != nil {
			return nil, fmt.Errorf("error creating otel exporter: %w", err)
		}
	default:
		return nil, fmt.Errorf("error creating otel meter: unknown mode %v", cfg.Mode)
	}

	ctrl := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(histogram.WithExplicitBoundaries(defaultHistogramBuckets)),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
		controller.WithExporter(exporter),
		controller.WithCollectPeriod(cfg.FlushInterval),
		controller.WithResource(resource.NewWithAttributes(instrumentationName, attribute.Key("metric.category").String("system"))),
	)
	if err = ctrl.Start(ctx); err != nil {
		return nil, fmt.Errorf("error starting push controller: %w", err)
	}

	global.SetMeterProvider(ctrl)
	return &OtelMeter{
		meter: global.Meter(instrumentationName),
		cfg:   cfg,
	}, nil
}

func (meter OtelMeter) DefaultHistogram(ctx context.Context, metricName string, value float64, fields util.ExtraFields) error {
	h, err := meter.meter.SyncFloat64().Histogram(metricName)
	if err != nil {
		return err
	}
	h.Record(ctx, value, append(fields.ToAttrs(), meter.defaultAttrs()...)...)
	return nil
}

func (meter OtelMeter) DefaultGauge(ctx context.Context, metricName string, value int64, fields util.ExtraFields) error {
	gauge, err := meter.meter.AsyncInt64().Gauge(metricName)
	if err != nil {
		return err
	}

	if err := meter.meter.RegisterCallback(
		[]instrument.Asynchronous{
			gauge,
		},
		func(ctx context.Context) {
			gauge.Observe(ctx, value, append(fields.ToAttrs(), meter.defaultAttrs()...)...)
		},
	); err != nil {
		return err
	}

	return nil
}

func (meter OtelMeter) DefaultCounter(ctx context.Context, metricName string, value int64, fields util.ExtraFields) error {
	counter, err := meter.meter.SyncInt64().Counter(metricName)
	if err != nil {
		return err
	}

	counter.Add(ctx, value, append(fields.ToAttrs(), meter.defaultAttrs()...)...)
	return nil
}

func (meter OtelMeter) defaultAttrs() []attribute.KeyValue {
	stackName := getStackName()
	defaultAttr := []attribute.KeyValue{
		attribute.Key("garden.app.name").String(meter.cfg.Service.Name),
		attribute.Key("garden.app.version").String(meter.cfg.Service.Version),
		attribute.Key("garden.stack").String(stackName),
	}

	if meter.cfg.DefaultFields != nil {
		for fieldName, value := range *meter.cfg.DefaultFields {
			defaultAttr = append(defaultAttr,
				attribute.Key(fieldName).String(value))
		}
	}

	return defaultAttr
}

func getStackName() string {
	stackName := os.Getenv("garden_STACK")
	if stackName == "" {
		stackName = "-"
	}
	return stackName
}
