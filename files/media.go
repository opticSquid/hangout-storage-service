package files

import (
	"github.com/knadh/koanf/v2"
	"hangout.com/core/storage-service/logger"
)

type mediaFile interface {
	processMedia(workerId int, cfg *koanf.Koanf, log logger.Log) error
}
