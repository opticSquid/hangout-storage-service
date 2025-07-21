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
		processDuration, _ = meter.Float64Histogram(metricName, metric.WithDescription(metricDesc), metric.WithExplicitBucketBoundaries(0, 30, 60, 90, 120, 150, 180, 210, 240, 270, 300, 330, 360, 390, 420, 450, 480, 510, 540, 570, 600))
	})
}
