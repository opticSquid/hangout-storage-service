package logger

import (
	"github.com/knadh/koanf/v2"
	"golang.org/x/net/context"
)

type Log interface {
	Debug(ctx context.Context, msg string, keysAndValues ...interface{})
	Info(ctx context.Context, msg string, keysAndValues ...interface{})
	Warn(ctx context.Context, msg string, keysAndValues ...interface{})
	Error(ctx context.Context, msg string, keysAndValues ...interface{})
}

func NewLogger(cfg *koanf.Koanf) Log {
	return NewZeroLogger(cfg)
}
