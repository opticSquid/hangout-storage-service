package kafka

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/knadh/koanf/v2"
	"hangout.com/core/storage-service/exceptions"
	"hangout.com/core/storage-service/files"
	"hangout.com/core/storage-service/logger"
)

// StartConsumer initializes and starts the Kafka consumer with the provided context
// This implements Sarama consumer group API
// Supports multi instance. All instances will join same consumer group.
// Single consumer per instance
func StartConsumer(eventChan chan<- *files.File, ctx context.Context, cfg *koanf.Koanf, log logger.Log) error {
	log.Debug("Configuring kafka client")
	consumerGroup, err := configureKafka(cfg)
	if err != nil {
		exceptions.KafkaConnectError("could not setup kafka connection", &err, log)
		return err
	}

	log.Info("kafka consumer group configured")
	go consume(eventChan, consumerGroup, ctx, cfg, log)
	return nil
}
func configureKafka(cfg *koanf.Koanf) (sarama.ConsumerGroup, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Version = sarama.DefaultVersion
	kafkaConfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	brokers := []string{cfg.String("kafka.url")}
	return sarama.NewConsumerGroup(brokers, cfg.String("kafka.groupId"), kafkaConfig)
}

func consume(eventChan chan<- *files.File, consumerGroup sarama.ConsumerGroup, ctx context.Context, cfg *koanf.Koanf, log logger.Log) {
	defer close(eventChan) // Close the channel when done
	defer consumerGroup.Close()

	handler := &ConsumerGroupHandler{Files: eventChan, log: log}
	for {
		select {
		case <-ctx.Done(): // Exit if the context is canceled
			log.Info("Context cancelled, stopping consumer")
			return
		default:
			if err := consumerGroup.Consume(ctx, []string{cfg.String("kafka.topic")}, handler); err != nil {
				exceptions.KafkaConsumerError("Error in consumer loop", &err, log)
				return
			}
		}
	}
}
