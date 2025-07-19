package logger

import (
	"github.com/knadh/koanf/v2"
	"golang.org/x/net/context"
)

type Log interface {
	Debug(ctx context.Context, msg string, keysAndValues ...any)
	Info(ctx context.Context, msg string, keysAndValues ...any)
	Warn(ctx context.Context, msg string, keysAndValues ...any)
	Error(ctx context.Context, msg string, keysAndValues ...any)
	With(keysAndValues ...any) Log
}

func NewLogger(cfg *koanf.Koanf) Log {
	return NewZeroLogger(cfg)
}
