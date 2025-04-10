package alog_test

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/arthurdotwork/alog"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := alog.Logger(
		alog.WithOutput(&buf),
		alog.WithLevel(slog.LevelInfo),
		alog.WithSource(true),
		alog.WithAttrs(slog.Attr{Key: "logger", Value: slog.StringValue("alog")}),
	)

	ctx := alog.Append(nil, slog.String("key", "value")) //nolint:staticcheck
	ctx = alog.Append(ctx, slog.Int64("count", 42))

	// act
	logger.InfoContext(ctx, "test message")

	// assert
	output := buf.String()
	require.Contains(t, output, `"message":"test message"`)
	require.Contains(t, output, `"severity":"info"`)
	require.Contains(t, output, `"logger":"alog"`)
	require.Contains(t, output, `"key":"value"`)
	require.Contains(t, output, `"count":42`)
}
