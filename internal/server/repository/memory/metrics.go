// Package memory provides persistent storage implementation that saves metrics data in memory
package memory

import (
	"context"
	"fmt"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/model"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/model/enum"
	"github.com/ruslanDantsov/osmetrics-server/internal/server/constants"
	"go.uber.org/zap"
	"sync"
)

// MemStorage реализация хранения метрик в памяти с потокобезопасным доступом.
type MemStorage struct {
	Mu      sync.RWMutex
	Storage map[string]*model.Metrics
	Log     zap.Logger
}

// NewMemStorage создает и возвращает новый экземпляр хранилища в памяти.
func NewMemStorage(log zap.Logger) *MemStorage {
	return &MemStorage{
		Storage: make(map[string]*model.Metrics),
		Log:     log,
	}
}

// HealthCheck проверяет состояние MemStorage.
func (s *MemStorage) HealthCheck(ctx context.Context) error {
	return nil
}

// Close закрывает MemStorage.
func (s *MemStorage) Close() {
	//For this type of storage we don't need implementation
}

// GetKnownMetrics возвращает список всех известных метрик
func (s *MemStorage) GetKnownMetrics(ctx context.Context) []string {
	metricNames := make([]string, 0, len(s.Storage))
	for name := range s.Storage {
		metricNames = append(metricNames, name)
	}
	return metricNames
}

// GetMetric возвращает метрику по заданному идентификатору.
func (s *MemStorage) GetMetric(ctx context.Context, metricID enum.MetricID) (*model.Metrics, bool) {
	if val, found := s.Storage[metricID.String()]; found {
		if val.MType == constants.CounterMetricType {
			s.Log.Info(fmt.Sprintf("Get metric name=%v type=%v delta=%v", val.ID, val.MType, *val.Delta))
		}

		if val.MType == constants.GaugeMetricType {
			s.Log.Info(fmt.Sprintf("Get metric name=%v type=%v value=%v", val.ID, val.MType, *val.Value))
		}
		return val, true
	}

	return nil, false
}

// SaveMetric сохраняет или обновляет одну метрику в хранилище.
func (s *MemStorage) SaveMetric(ctx context.Context, metric *model.Metrics) (*model.Metrics, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	key := metric.ID.String()
	existing, found := s.Storage[key]

	if metric.MType == constants.CounterMetricType {
		if found && existing.Delta != nil && metric.Delta != nil {
			*existing.Delta += *metric.Delta
			s.Log.Info(fmt.Sprintf("UPDATE counter_metric name=%v delta=%v", metric.ID, *existing.Delta))
			return existing, nil
		}
	}

	s.Storage[key] = metric

	if metric.MType == constants.CounterMetricType {
		s.Log.Info(fmt.Sprintf("SAVE %v metric id=%v delta=%v", metric.MType, metric.ID, *metric.Delta))
	}

	if metric.MType == constants.GaugeMetricType {
		s.Log.Info(fmt.Sprintf("SAVE %v metric id=%v value=%v", metric.MType, metric.ID, *metric.Value))
	}

	return metric, nil
}

// SaveAllMetrics сохраняет список метрик в хранилище.
func (s *MemStorage) SaveAllMetrics(ctx context.Context, metricList model.MetricsList) (model.MetricsList, error) {
	var savedMetrics model.MetricsList

	for _, metric := range metricList {
		savedMetric, err := s.SaveMetric(ctx, &metric)
		if err != nil {
			return nil, fmt.Errorf("failed to save metric %v: %w", metric.ID, err)
		}
		savedMetrics = append(savedMetrics, *savedMetric)
	}

	return savedMetrics, nil
}
