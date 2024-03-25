package logging

import (
	"io"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

type LoggerConfig struct {
	Level  slog.Level
	Pretty bool
	w      io.Writer
}

type Option func(*LoggerConfig)

func WithLevel(level string) Option {
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		lvl = slog.LevelError
	}

	return func(c *LoggerConfig) {
		c.Level = lvl
	}
}

func WithPretty(pretty bool) Option {
	return func(c *LoggerConfig) {
		c.Pretty = pretty
	}
}

func WithWriter(w io.Writer) Option {
	return func(c *LoggerConfig) {
		c.w = w
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

	if options.w == nil {
		options.w = os.Stdout
	}

	if options.Pretty {
		handler = tint.NewHandler(options.w, &tint.Options{
			AddSource: true,
			Level:     options.Level,
		})
	} else {
		handler = slog.NewJSONHandler(options.w, &slog.HandlerOptions{
			AddSource: true,
			Level:     options.Level,
		})
	}

	// create a logger
	logger := slog.New(handler)

	slog.SetDefault(logger)

	return logger
}
