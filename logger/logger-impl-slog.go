package logger

import (
	"log/slog"
	"os"

	"hangout.com/core/storage-service/config"
)

type SlogLogger struct {
	log *slog.Logger
}

func NewSlogLogger(cfg *config.Config) Log {
	var logLevel slog.Level
	switch cfg.Log.Level {
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
	return &SlogLogger{log: sl}
}

func (sl *SlogLogger) Debug(message string, keysAndValues ...interface{}) {
	sl.log.Debug(message, keysAndValues...)
}

func (sl *SlogLogger) Info(message string, keysAndValues ...interface{}) {
	sl.log.Info(message, keysAndValues...)
}

func (sl *SlogLogger) Warn(message string, keysAndValues ...interface{}) {
	sl.log.Warn(message, keysAndValues...)
}

func (sl *SlogLogger) Error(message string, keysAndValues ...interface{}) {
	sl.log.Error(message, keysAndValues...)
}