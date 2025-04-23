package service

import (
	"github.com/go-resty/resty/v2"
	"github.com/ruslanDantsov/osmetrics-server/internal/config"
	"github.com/ruslanDantsov/osmetrics-server/internal/model/enum"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

type MockRestClient struct{}

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
func setupTest() (*zap.Logger, *MockRestClient, *config.AgentConfig) {
	logger, _ := zap.NewDevelopment()
	return logger, &MockRestClient{}, NewMockAgentConfig()
}

func TestAppendMetric(t *testing.T) {
	logger, client, cfg := setupTest()
	defer logger.Sync()

	ms := NewMetricService(logger, client, cfg)

	ms.appendMetric(enum.Alloc, 123.45)

	val, exists := ms.Metrics[enum.Alloc]
	assert.True(t, exists)
	assert.Equal(t, 123.45, val)
}

func TestAggregateMetric(t *testing.T) {
	logger, client, cfg := setupTest()
	defer logger.Sync()

	ms := NewMetricService(logger, client, cfg)
	ms.aggregateMetric(enum.PollCount, 5)
	ms.aggregateMetric(enum.PollCount, 3)

	val, exists := ms.Metrics[enum.PollCount]
	assert.True(t, exists)
	assert.Equal(t, int64(8), val)
}

func TestCollectMetrics(t *testing.T) {
	logger, client, cfg := setupTest()
	defer logger.Sync()

	ms := NewMetricService(logger, client, cfg)

	ms.CollectMetrics()

	_, exists := ms.Metrics[enum.Alloc]
	assert.True(t, exists)
	_, exists = ms.Metrics[enum.RandomValue]
	assert.True(t, exists)
}
