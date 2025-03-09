package alog

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"

	slogctx "github.com/veqryn/slog-context"
)

type LoggerOptions struct {
	Output io.Writer
	Level  slog.Level
	Source bool

	Attrs []slog.Attr
}

type LoggerOptionsFunc func(*LoggerOptions)

func WithOutput(w io.Writer) LoggerOptionsFunc {
	return func(o *LoggerOptions) {
		o.Output = w
	}
}

func WithLevel(l slog.Level) LoggerOptionsFunc {
	return func(o *LoggerOptions) {
		o.Level = l
	}
}

func WithSource(b bool) LoggerOptionsFunc {
	return func(o *LoggerOptions) {
		o.Source = b
	}
}

func WithAttrs(attrs ...slog.Attr) LoggerOptionsFunc {
	return func(o *LoggerOptions) {
		o.Attrs = attrs
	}
}

func Logger(opts ...LoggerOptionsFunc) *slog.Logger {
	options := LoggerOptions{
		Output: os.Stdout,
		Level:  slog.LevelInfo,
		Source: true,
		Attrs:  nil,
	}

	for _, opt := range opts {
		opt(&options)
	}

	jsonLogHandler := slog.NewJSONHandler(options.Output, &slog.HandlerOptions{
		AddSource: true,
		Level:     options.Level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "msg" {
				a.Key = "message"
			}

			if a.Key == "level" {
				a.Key = "severity"
				a.Value = slog.StringValue(strings.ToLower(a.Value.String()))
			}

			return a
		},
	}).WithAttrs(options.Attrs)

	logger := slog.New(slogctx.NewHandler(jsonLogHandler, nil))
	return logger
}

func Append(ctx context.Context, key string, value any) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx = slogctx.Append(ctx, key, value)
	return ctx
}
