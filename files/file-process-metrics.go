package files

import (
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var (
	processDuration metric.Float64Histogram
	initOnce        sync.Once
)

func initFileProcessMetrics(metricName string, metricDesc string) {
	// Initialize file process metrics here
	initOnce.Do(func() {
		meter := otel.GetMeterProvider().Meter("hangout.storage.file")
		processDuration, _ = meter.Float64Histogram(metricName, metric.WithDescription(metricDesc), metric.WithExplicitBucketBoundaries(5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60))
	})
}
