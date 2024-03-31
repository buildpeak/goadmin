package api

import (
	"testing"
)

func TestStartOtel(t *testing.T) {
	t.Parallel()

	type args struct {
		serviceInfo         ServiceInfo
		observabilityConfig ObservabilityConfig
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				serviceInfo: ServiceInfo{
					Name:    "service",
					Version: "1.0.0",
					Env:     "test",
				},
				observabilityConfig: ObservabilityConfig{
					Collector: Collector{
						Host:               "localhost",
						Port:               4317,
						Headers:            []Header{{Key: "key", Value: "value"}},
						IsInsecure:         true,
						WithMetricsEnabled: true,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := StartOtel(tt.args.serviceInfo, tt.args.observabilityConfig)

			if (err != nil) != tt.wantErr {
				t.Errorf("StartOtel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
