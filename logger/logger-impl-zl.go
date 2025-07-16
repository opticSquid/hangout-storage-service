package logger

import (
	"os"

	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
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
	zl := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &zeroLogger{log: zl}
}

func (zl *zeroLogger) Debug(msg string, keysAndValues ...interface{}) {
	zl.log.Debug().Fields(keysAndValues).Msg(msg)
}

func (zl *zeroLogger) Info(msg string, keysAndValues ...interface{}) {
	zl.log.Info().Fields(keysAndValues).Msg(msg)
}

func (zl *zeroLogger) Warn(msg string, keysAndValues ...interface{}) {
	zl.log.Warn().Fields(keysAndValues).Msg(msg)
}

func (zl *zeroLogger) Error(msg string, keysAndValues ...interface{}) {
	zl.log.Error().Fields(keysAndValues).Msg(msg)
}
