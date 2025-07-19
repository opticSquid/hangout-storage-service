package logger

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"

	"github.com/knadh/koanf/v2"
)

type slogLogger struct {
	log *slog.Logger
}

func NewSlogLogger(cfg *koanf.Koanf) Log {
	var logLevel slog.Level
	switch cfg.String("log.level") {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}
	sl := otelslog.NewLogger(cfg.String("application.name"))
	slog.SetLogLoggerLevel(logLevel)
	slog.SetDefault(sl)
	return &slogLogger{log: sl}
}

func (sl *slogLogger) Debug(ctx context.Context, message string, keysAndValues ...any) {
	sl.log.DebugContext(ctx, message, keysAndValues...)
}

func (sl *slogLogger) Info(ctx context.Context, message string, keysAndValues ...any) {
	sl.log.InfoContext(ctx, message, keysAndValues...)
}

func (sl *slogLogger) Warn(ctx context.Context, message string, keysAndValues ...any) {
	sl.log.WarnContext(ctx, message, keysAndValues...)
}

func (sl *slogLogger) Error(ctx context.Context, message string, keysAndValues ...any) {
	sl.log.ErrorContext(ctx, message, keysAndValues...)
}
