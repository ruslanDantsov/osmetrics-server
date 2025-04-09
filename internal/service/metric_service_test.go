package service

import (
	"github.com/go-resty/resty/v2"
	"github.com/ruslanDantsov/osmetrics-server/internal/config"
	"github.com/ruslanDantsov/osmetrics-server/internal/model/enum/metric"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockLogger struct{}

func (m MockLogger) Info(msg string)  {}
func (m MockLogger) Error(msg string) {}

type MockRestClient struct{}

type MockConfig struct{}

func NewMockAgentConfig() *config.AgentConfig {
	return &config.AgentConfig{
		Address:        "localhost:8080",
		ReportInterval: 10,
		PollInterval:   2,
	}
}

func (mrc *MockRestClient) R() *resty.Request {
	return resty.New().R().
		SetDoNotParseResponse(true).
		SetBody("mock").
		SetHeader("Content-Type", "text/plain")
}

func TestAppendMetric(t *testing.T) {
	ms := NewMetricService(MockLogger{}, &MockRestClient{}, NewMockAgentConfig())
	ms.appendMetric(metric.Alloc, 123.45)

	val, exists := ms.Metrics[metric.Alloc]
	assert.True(t, exists)
	assert.Equal(t, 123.45, val)
}

func TestAggregateMetric(t *testing.T) {
	ms := NewMetricService(MockLogger{}, &MockRestClient{}, NewMockAgentConfig())
	ms.aggregateMetric(metric.PollCount, 5)
	ms.aggregateMetric(metric.PollCount, 3)

	val, exists := ms.Metrics[metric.PollCount]
	assert.True(t, exists)
	assert.Equal(t, int64(8), val)
}

func TestCollectMetrics(t *testing.T) {
	ms := NewMetricService(MockLogger{}, &MockRestClient{}, NewMockAgentConfig())
	ms.CollectMetrics()

	_, exists := ms.Metrics[metric.Alloc]
	assert.True(t, exists)
	_, exists = ms.Metrics[metric.RandomValue]
	assert.True(t, exists)
}
