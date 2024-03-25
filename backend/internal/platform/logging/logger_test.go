package logging

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"
)

func TestWithLevel(t *testing.T) {
	t.Parallel()

	type args struct {
		level string
	}

	tests := []struct {
		name string
		args args
		cfg  *LoggerConfig
	}{
		{
			name: "TestWithLevel",
			args: args{
				level: "info",
			},
			cfg: &LoggerConfig{},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			WithLevel(tt.args.level)(tt.cfg)

			if tt.cfg.Level != slog.LevelInfo {
				t.Errorf("WithLevel()() = %v, want %v", tt.cfg.Level, slog.LevelInfo)
			}
		})
	}
}

func TestWithPretty(t *testing.T) {
	t.Parallel()

	type args struct {
		pretty bool
	}

	tests := []struct {
		name string
		args args
		cfg  *LoggerConfig
	}{
		{
			name: "TestWithPretty",
			args: args{
				pretty: true,
			},
			cfg: &LoggerConfig{},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			WithPretty(tt.args.pretty)(tt.cfg)

			if tt.cfg.Pretty != true {
				t.Errorf("WithPretty()() = %v, want %v", tt.cfg.Pretty, true)
			}
		})
	}
}

func TestWithWriter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		buf  *bytes.Buffer
		cfg  *LoggerConfig
	}{
		{
			name: "TestWithWriter",
			buf:  bytes.NewBuffer(nil),
			cfg:  &LoggerConfig{},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			WithWriter(tt.buf)(tt.cfg)

			if tt.cfg.w != tt.buf {
				t.Errorf("WithWriter()() = %v, want %v", tt.cfg.w, tt.buf)
			}
		})
	}
}

func TestNewLogger(t *testing.T) {
	t.Parallel()

	type args struct {
		opts []Option
	}

	tests := []struct {
		name string
		args args
		buf  *bytes.Buffer
		msg  string
		want string
	}{
		{
			name: "TestNewLogger",
			args: args{
				opts: []Option{
					WithLevel("info"),
				},
			},
			buf:  bytes.NewBuffer([]byte{}),
			msg:  "test",
			want: "test",
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.args.opts = append(tt.args.opts, WithPretty(false), WithWriter(tt.buf))

			logger := NewLogger(tt.args.opts...)

			if logger == nil {
				t.Errorf("NewLogger() = %v, want %v", logger, nil)
			}

			logger.Info(tt.msg)

			got := make(map[string]interface{})
			if err := json.Unmarshal(tt.buf.Bytes(), &got); err != nil {
				t.Errorf("json.Unmarshal() = %v", err)
			}

			t.Log(got)

			if got["msg"] != tt.want {
				t.Errorf("NewLogger() = %v, want %v", got["msg"], tt.want)
			}

			if got["level"] != "INFO" {
				t.Errorf("NewLogger() = %v, want %v", got["level"], "INFO")
			}

			if got["source"] == nil {
				t.Errorf("NewLogger() = %v, want %v", got["source"], nil)
			}
		})
	}
}
