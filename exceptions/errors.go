package exceptions

import (
	"context"
	"os"

	"hangout.com/core/storage-service/logger"
)

func ProcessError(ctx context.Context, msg string, err *error, log logger.Log) {
	log.Error(ctx, msg, "error", err)
	os.Exit(2)
}

func KafkaConnectError(ctx context.Context, msg string, err *error, log logger.Log) {
	log.Error(ctx, msg, "error", err)
	os.Exit(3)
}

func KafkaConsumerError(ctx context.Context, msg string, err *error, log logger.Log) {
	log.Error(ctx, msg, "error", err)
	os.Exit(4)
}

func DbConnectionError(ctx context.Context, msg string, err *error, log logger.Log) {
	log.Error(ctx, msg, "error", err)
	os.Exit(5)
}
