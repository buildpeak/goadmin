package api

import (
	"fmt"

	"github.com/honeycombio/otel-config-go/otelconfig"
)

type Header struct {
	Key   string
	Value string
}

type ObservabilityConfig struct {
	Collector struct {
		Host               string
		Port               int
		Headers            []Header
		IsInsecure         bool
		WithMetricsEnabled bool
	}
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
