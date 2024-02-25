package logging

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

type LoggerConfig struct {
	Level  slog.Level
	Pretty bool
}

type Option func(*LoggerConfig)

func WithLevel(level string) Option {
	var l slog.Level
	if err := l.UnmarshalText([]byte(level)); err != nil {
		l = slog.LevelError
	}
	return func(c *LoggerConfig) {
		c.Level = l
	}
}

func WithPretty(pretty bool) Option {
	return func(c *LoggerConfig) {
		c.Pretty = pretty
	}
}

func NewLogger(opts ...Option) *slog.Logger {
	options := &LoggerConfig{
		Level:  slog.LevelInfo,
		Pretty: false,
	}

	for _, opt := range opts {
		opt(options)
	}

	// create a handler
	var handler slog.Handler

	if options.Pretty {
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			AddSource: true,
			Level:     options.Level,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     options.Level,
		})
	}

	// create a logger
	logger := slog.New(handler)

	slog.SetDefault(logger)

	return logger
}
