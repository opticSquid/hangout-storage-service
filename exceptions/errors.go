package exceptions

import (
	"os"

	"hangout.com/core/storage-service/logger"
)

func ProcessError(msg string, err *error, log logger.Log) {
	log.Error(msg, "error", err)
	os.Exit(2)
}

func KafkaConnectError(msg string, err *error, log logger.Log) {
	log.Error(msg, "error", err)
	os.Exit(3)
}

func KafkaConsumerError(msg string, err *error, log logger.Log) {
	log.Error(msg, "error", err)
	os.Exit(4)
}

func DbConnectionError(msg string, err *error, log logger.Log) {
	log.Error(msg, "error", err)
	os.Exit(5)
}
