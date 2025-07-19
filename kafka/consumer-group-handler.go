package kafka

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"hangout.com/core/storage-service/files"
	"hangout.com/core/storage-service/logger"
)

// ConsumerGroupHandler implements sarama.ConsumerGroupHandler
type ConsumerGroupHandler struct {
	Files chan<- *files.File
	ctx   context.Context
	log   logger.Log
}

// Setup runs at the beginning of a new session, before ConsumeClaim
func (cgh *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	cgh.log.Info(cgh.ctx, "Consumer group session setup completed")
	return nil
}

// Cleanup runs at the end of a session, once all ConsumeClaim goroutines have exited
func (cgh *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	cgh.log.Info(cgh.ctx, "Consumer group session cleanup completed")
	return nil
}

// Helper type to adapt Sarama headers to OpenTelemetry TextMapCarrier
type kafkaHeaderCarrier []*sarama.RecordHeader

func (c kafkaHeaderCarrier) Get(key string) string {
	for _, h := range c {
		if string(h.Key) == key {
			return string(h.Value)
		}
	}
	return ""
}
func (c kafkaHeaderCarrier) Set(key, value string) {}
func (c kafkaHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for _, h := range c {
		keys = append(keys, string(h.Key))
	}
	return keys
}

// ConsumeClaim starts a consumer loop for each partition assigned to this handler
func (cgh *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	tr := otel.Tracer("hangout.storage.kafka")
	propagator := otel.GetTextMapPropagator()
	for message := range claim.Messages() {
		carrier := kafkaHeaderCarrier(message.Headers)
		parentCtx := propagator.Extract(cgh.ctx, carrier)
		ctx, span := tr.Start(parentCtx, "ConsumeKafkaMessage")
		span.SetAttributes(
			attribute.String("type", "consumer"),
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.destination", message.Topic),
			attribute.Int64("messaging.kafka.partition", int64(message.Partition)),
			attribute.Int64("messaging.kafka.offset", message.Offset),
		)
		var body eventBody
		err := json.Unmarshal(message.Value, &body)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			cgh.log.Error(ctx, "error in unmarshalling", err)
		} else {
			event := files.File{Context: ctx, Filename: body.Filename, ContentType: body.ContentType, UserId: body.UserId}
			cgh.log.Debug(ctx, "File Upload event occured",
				"Topic", message.Topic,
				"Partition", message.Partition,
				"Offset", message.Offset,
				"Header", message.Headers,
				"Value", string(message.Value),
			)
			select {
			case cgh.Files <- &event:
				session.MarkMessage(message, "")
				span.End()
			default:
				cgh.log.Warn(ctx, "File channel is full, unable to process event",
					"FileName", event.Filename,
					"ContentType", event.ContentType,
					"Partition", message.Partition,
					"Offset", message.Offset,
					"Header", message.Headers,
					"Value", string(message.Value),
				)
				span.End()
			}
		}

	}
	return nil
}
