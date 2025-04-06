package service

import (
	"github.com/ruslanDantsov/osmetrics-server/internal/model/enum/metric"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockLogger struct{}

func (m MockLogger) Info(msg string)  {}
func (m MockLogger) Error(msg string) {}

func TestAppendMetric(t *testing.T) {
	ms := NewMetricService(MockLogger{})
	ms.appendMetric(metric.Alloc, 123.45)

	val, exists := ms.Metrics[metric.Alloc]
	assert.True(t, exists)
	assert.Equal(t, 123.45, val)
}

func TestAggregateMetric(t *testing.T) {
	ms := NewMetricService(MockLogger{})
	ms.aggregateMetric(metric.PollCount, 5)
	ms.aggregateMetric(metric.PollCount, 3)

	val, exists := ms.Metrics[metric.PollCount]
	assert.True(t, exists)
	assert.Equal(t, int64(8), val)
}

func TestCollectMetrics(t *testing.T) {
	ms := NewMetricService(MockLogger{})
	ms.CollectMetrics()

	_, exists := ms.Metrics[metric.Alloc]
	assert.True(t, exists)
	_, exists = ms.Metrics[metric.RandomValue]
	assert.True(t, exists)
}
