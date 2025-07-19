package logger

import (
	"context"
	"os"

	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/bridges/otelzerolog"
)

type zeroLogger struct {
	log zerolog.Logger
}

func NewZeroLogger(cfg *koanf.Koanf) Log {
	switch cfg.String("log.level") {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	hook := otelzerolog.NewHook(cfg.String("application.name"))
	zl := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &zeroLogger{log: zl.Hook(hook)}
}

func (zl *zeroLogger) Debug(ctx context.Context, msg string, keysAndValues ...any) {
	zl.log.Debug().Ctx(ctx).Fields(keysAndValues).Msg(msg)
}

func (zl *zeroLogger) Info(ctx context.Context, msg string, keysAndValues ...any) {
	zl.log.Info().Ctx(ctx).Fields(keysAndValues).Msg(msg)
}

func (zl *zeroLogger) Warn(ctx context.Context, msg string, keysAndValues ...any) {
	zl.log.Warn().Ctx(ctx).Fields(keysAndValues).Msg(msg)
}

func (zl *zeroLogger) Error(ctx context.Context, msg string, keysAndValues ...any) {
	zl.log.Error().Ctx(ctx).Fields(keysAndValues).Msg(msg)
}
