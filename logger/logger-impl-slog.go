package logger

import (
	"log/slog"
	"os"

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
	sl := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(sl)
	return &slogLogger{log: sl}
}

func (sl *slogLogger) Debug(message string, keysAndValues ...interface{}) {
	sl.log.Debug(message, keysAndValues...)
}

func (sl *slogLogger) Info(message string, keysAndValues ...interface{}) {
	sl.log.Info(message, keysAndValues...)
}

func (sl *slogLogger) Warn(message string, keysAndValues ...interface{}) {
	sl.log.Warn(message, keysAndValues...)
}

func (sl *slogLogger) Error(message string, keysAndValues ...interface{}) {
	sl.log.Error(message, keysAndValues...)
}
