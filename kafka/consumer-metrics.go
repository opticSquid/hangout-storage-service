package kafka

import (
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var (
	consumedEventsCounter metric.Int64Counter
	initMetricsOnce       sync.Once
)

func initConsumerMetrics() {
	initMetricsOnce.Do(func() {
		meter := otel.GetMeterProvider().Meter("hangout.storage.kafka")
		consumedEventsCounter, _ = meter.Int64Counter(
			"kafka_consumer_events_consumed_total",
			metric.WithDescription("Total number of Kafka events consumed"),
		)
	})
}
