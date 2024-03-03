package api

import (
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		want    *Config
		wantErr bool
	}{
		{
			name: "success",
			want: &Config{
				Env:         "development",
				DatabaseURL: "postgres://localhost:5432/mydb",
				Log: struct {
					Level  string `json:"level"`
					Pretty bool   `json:"pretty"`
				}{
					Level:  "info",
					Pretty: true,
				},
				APIServer: APIServerConfig{
					Port: 8080,
					Auth: struct {
						JWTSecret string `json:"jwt_secret"`
					}{
						JWTSecret: "secret",
					},
				},
				Observability: ObservabilityConfig{
					Collector: Collector{
						Host:               "localhost",
						Port:               13133,
						Headers:            []Header{{Key: "key", Value: "value"}},
						IsInsecure:         true,
						WithMetricsEnabled: true,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
