package metrics

import (
	"fmt"

	"github.com/garden/observability-commons/config"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
)

const (
	prodEndpoint  = "otel-collector.garden.internal"
	devEndpoint   = "localhost:4317"
	debugEndpoint = "localhost:4317"
)

func newClient(cfg config.Config) otlpmetric.Client {
	return otlpmetricgrpc.NewClient(
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(endpoint(cfg.Mode, cfg.Port)),
		otlpmetricgrpc.WithTimeout(cfg.Timeout),
	)
}

func endpoint(mode config.Mode, port string) string {
	switch mode {
	case config.Debug:
		return debugEndpoint
	case config.Development:
		return devEndpoint
	case config.Production:
		return fmt.Sprintf("%s:%s", prodEndpoint, port)
	default:
		return ""
	}
}
