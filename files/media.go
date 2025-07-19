package files

import (
	"context"

	"github.com/knadh/koanf/v2"
	"hangout.com/core/storage-service/logger"
)

type mediaFile interface {
	processMedia(workerId int, ctx context.Context, cfg *koanf.Koanf, log logger.Log) error
}
