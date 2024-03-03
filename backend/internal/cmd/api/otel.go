package api

import (
	"fmt"

	"github.com/honeycombio/otel-config-go/otelconfig"
)

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Collector struct {
	Host               string   `json:"host"`
	Port               int      `json:"port"`
	Headers            []Header `json:"headers"`
	IsInsecure         bool     `json:"is_insecure"`
	WithMetricsEnabled bool     `json:"with_metrics_enabled"`
}

type ObservabilityConfig struct {
	Collector Collector `json:"collector"`
}

type ServiceInfo struct {
	Name    string
	Version string
	Env     string
}

func StartOtel(serviceInfo ServiceInfo, observabilityConfig ObservabilityConfig) (func(), error) {
	headers := make(map[string]string, len(observabilityConfig.Collector.Headers))
	for _, header := range observabilityConfig.Collector.Headers {
		headers[header.Key] = header.Value
	}

	shutdown, err := otelconfig.ConfigureOpenTelemetry(
		otelconfig.WithServiceName(fmt.Sprintf("%s-%s", serviceInfo.Name, serviceInfo.Env)),
		otelconfig.WithServiceVersion(serviceInfo.Version),
		otelconfig.WithHeaders(headers),
		otelconfig.WithExporterEndpoint(
			fmt.Sprintf("%s:%d", observabilityConfig.Collector.Host, observabilityConfig.Collector.Port),
		),
		otelconfig.WithExporterInsecure(observabilityConfig.Collector.IsInsecure),
		otelconfig.WithMetricsEnabled(observabilityConfig.Collector.WithMetricsEnabled),
	)
	if err != nil {
		return shutdown, fmt.Errorf("failed to start otel: %w", err)
	}

	return shutdown, nil
}
