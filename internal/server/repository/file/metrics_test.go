package file

import (
	"context"
	"errors"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/model"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/model/enum"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"os"
	"testing"
)

type MockMemoryStorager struct {
	mock.Mock
	Storage map[string]*model.Metrics
}

func (m *MockMemoryStorager) GetKnownMetrics(ctx context.Context) []string {
	args := m.Called(ctx)
	return args.Get(0).([]string)
}

func (m *MockMemoryStorager) GetMetric(ctx context.Context, metricID enum.MetricID) (*model.Metrics, bool) {
	args := m.Called(ctx, metricID)
	return args.Get(0).(*model.Metrics), args.Bool(1)
}

func (m *MockMemoryStorager) SaveMetric(ctx context.Context, metric *model.Metrics) (*model.Metrics, error) {
	args := m.Called(ctx, metric)
	return args.Get(0).(*model.Metrics), args.Error(1)
}

func (m *MockMemoryStorager) SaveAllMetrics(ctx context.Context, metricList model.MetricsList) (model.MetricsList, error) {
	args := m.Called(ctx, metricList)
	return args.Get(0).(model.MetricsList), args.Error(1)
}

func (m *MockMemoryStorager) HealthCheck(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMemoryStorager) Close() {
	m.Called()
}

func TestPersistentStorage_SaveMetric(t *testing.T) {
	mockStorage := new(MockMemoryStorager)
	logger := zap.NewNop()
	filePath := "./test_save_metric.json"
	defer os.Remove(filePath)

	ctx := context.Background()
	metric := &model.Metrics{ID: "Alloc", MType: "gauge", Value: floatPointer(123.4)}
	mockStorage.On("SaveMetric", ctx, metric).Return(metric, nil)

	ps := NewPersistentStorage(mockStorage, filePath, 0, *logger, false)

	result, err := ps.SaveMetric(ctx, metric)
	assert.NoError(t, err)
	assert.Equal(t, metric, result)

	mockStorage.AssertExpectations(t)
}

func TestPersistentStorage_SaveMetric_Error(t *testing.T) {
	mockStorage := new(MockMemoryStorager)
	logger := zap.NewNop()
	filePath := "./test_save_metric_err.json"
	defer os.Remove(filePath)

	ctx := context.Background()
	metric := &model.Metrics{ID: "Heap", MType: "gauge", Value: floatPointer(11.1)}
	mockErr := errors.New("save failed")
	mockStorage.On("SaveMetric", ctx, metric).Return(metric, mockErr)

	ps := NewPersistentStorage(mockStorage, filePath, 0, *logger, false)

	result, err := ps.SaveMetric(ctx, metric)
	assert.Nil(t, result)
	assert.Equal(t, mockErr, err)

	mockStorage.AssertExpectations(t)
}

func TestPersistentStorage_SaveAllMetrics(t *testing.T) {
	mockStorage := new(MockMemoryStorager)
	logger := zap.NewNop()
	filePath := "./test_save_all.json"
	defer os.Remove(filePath)

	ctx := context.Background()
	metrics := model.MetricsList{
		{ID: "CPU", MType: "gauge", Value: floatPointer(55)},
	}
	mockStorage.On("SaveAllMetrics", ctx, metrics).Return(metrics, nil)

	ps := NewPersistentStorage(mockStorage, filePath, 0, *logger, false)

	result, err := ps.SaveAllMetrics(ctx, metrics)
	assert.NoError(t, err)
	assert.Equal(t, metrics, result)

	mockStorage.AssertExpectations(t)
}

func TestPersistentStorage_SaveAllMetrics_Error(t *testing.T) {
	mockStorage := new(MockMemoryStorager)
	logger := zap.NewNop()
	filePath := "./test_save_all_err.json"
	defer os.Remove(filePath)

	ctx := context.Background()
	metrics := model.MetricsList{
		{ID: "RAM", MType: "gauge", Value: floatPointer(88)},
	}
	mockErr := errors.New("bulk save failed")
	mockStorage.On("SaveAllMetrics", ctx, metrics).Return(metrics, mockErr)

	ps := NewPersistentStorage(mockStorage, filePath, 0, *logger, false)

	result, err := ps.SaveAllMetrics(ctx, metrics)
	assert.Nil(t, result)
	assert.Equal(t, mockErr, err)

	mockStorage.AssertExpectations(t)
}

func TestPersistentStorage_GetMetric(t *testing.T) {
	mockStorage := new(MockMemoryStorager)
	logger := zap.NewNop()

	ctx := context.Background()
	metricID := enum.MetricID("GC")
	expected := &model.Metrics{ID: "GC", MType: "gauge", Value: floatPointer(12)}
	mockStorage.On("GetMetric", ctx, metricID).Return(expected, true)

	ps := NewPersistentStorage(mockStorage, "", 0, *logger, false)

	result, ok := ps.GetMetric(ctx, metricID)
	assert.True(t, ok)
	assert.Equal(t, expected, result)

	mockStorage.AssertExpectations(t)
}

func TestPersistentStorage_GetKnownMetrics(t *testing.T) {
	mockStorage := new(MockMemoryStorager)
	logger := zap.NewNop()

	ctx := context.Background()
	expected := []string{"Alloc", "Heap"}
	mockStorage.On("GetKnownMetrics", ctx).Return(expected)

	ps := NewPersistentStorage(mockStorage, "", 0, *logger, false)

	result := ps.GetKnownMetrics(ctx)
	assert.Equal(t, expected, result)

	mockStorage.AssertExpectations(t)
}

func floatPointer(v float64) *float64 {
	return &v
}
