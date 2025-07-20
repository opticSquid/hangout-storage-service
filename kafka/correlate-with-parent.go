package kafka

import "github.com/IBM/sarama"

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
