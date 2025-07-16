package kafka

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"hangout.com/core/storage-service/files"
	"hangout.com/core/storage-service/logger"
)

// ConsumerGroupHandler implements sarama.ConsumerGroupHandler
type ConsumerGroupHandler struct {
	Files chan<- *files.File
	log   logger.Log
}

// Setup runs at the beginning of a new session, before ConsumeClaim
func (cgh *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	cgh.log.Info("Consumer group session setup completed")
	return nil
}

// Cleanup runs at the end of a session, once all ConsumeClaim goroutines have exited
func (cgh *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	cgh.log.Info("Consumer group session cleanup completed")
	return nil
}

// ConsumeClaim starts a consumer loop for each partition assigned to this handler
func (cgh *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var body eventBody
		err := json.Unmarshal(message.Value, &body)
		if err != nil {
			cgh.log.Error("error in unmarshalling", err)
		} else {
			event := files.File{Filename: body.Filename, ContentType: body.ContentType, UserId: body.UserId}
			cgh.log.Debug("File Upload event occured",
				"Topic", message.Topic,
				"Partition", message.Partition,
				"Offset", message.Offset,
				"Value", string(message.Value),
			)
			// Send event to Files channel with non-blocking check
			select {
			case cgh.Files <- &event:
				session.MarkMessage(message, "")
			default:
				// Log a warning if the channel is full or has no active consumers
				cgh.log.Warn("File channel is full, unable to process event",
					"FileName", event.Filename,
					"ContentType", event.ContentType,
					"Partition", message.Partition,
					"Offset", message.Offset,
				)
			}
		}

	}
	return nil
}
