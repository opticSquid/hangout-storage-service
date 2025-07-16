package logger

import (
	"github.com/knadh/koanf/v2"
)

type Log interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

func NewLogger(cfg *koanf.Koanf) Log {
	if cfg.String("log.backend") == "slog" {
		return NewSlogLogger(cfg)
	} else {
		return NewZeroLogger(cfg)
	}
}
