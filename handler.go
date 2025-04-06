package alog

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
)

// ContextKey is the type used for context keys in the alog package
type ContextKey string

// LogAttrKey is the context key used to store log attributes
const LogAttrKey ContextKey = "alog_attributes"

// LoggerOptions holds configuration for creating a new logger
type LoggerOptions struct {
	Output io.Writer
	Level  slog.Level
	Source bool
	Attrs  []slog.Attr
}

// LoggerOptionsFunc is a function that modifies LoggerOptions
type LoggerOptionsFunc func(*LoggerOptions)

// WithOutput sets the output writer for the logger
func WithOutput(w io.Writer) LoggerOptionsFunc {
	return func(o *LoggerOptions) {
		o.Output = w
	}
}

// WithLevel sets the minimum log level
func WithLevel(l slog.Level) LoggerOptionsFunc {
	return func(o *LoggerOptions) {
		o.Level = l
	}
}

// WithSource enables or disables source code location in logs
func WithSource(b bool) LoggerOptionsFunc {
	return func(o *LoggerOptions) {
		o.Source = b
	}
}

// WithAttrs adds default attributes to the logger
func WithAttrs(attrs ...slog.Attr) LoggerOptionsFunc {
	return func(o *LoggerOptions) {
		o.Attrs = attrs
	}
}

// Logger creates a new slog.Logger with the specified options
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
		AddSource: options.Source,
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

	// Wrap the JSON handler with our context-aware handler
	logger := slog.New(&contextHandler{Handler: jsonLogHandler})
	return logger
}

// contextHandler is a slog.Handler that extracts attributes from the context
type contextHandler struct {
	slog.Handler
}

// Handle implements slog.Handler.Handle by extracting attributes from ctx
func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if ctx != nil {
		if attrs, ok := ctx.Value(LogAttrKey).([]slog.Attr); ok && len(attrs) > 0 {
			for _, attr := range attrs {
				r.AddAttrs(attr)
			}
		}
	}

	return h.Handler.Handle(ctx, r)
}

// Append adds a key-value pair to the context for logging
func Append(ctx context.Context, attr slog.Attr) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	var attrs []slog.Attr
	if existingAttrs, ok := ctx.Value(LogAttrKey).([]slog.Attr); ok {
		attrs = existingAttrs
	}
	attrs = append(attrs, attr)

	return context.WithValue(ctx, LogAttrKey, attrs)
}
