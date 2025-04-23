package repository

import (
	"fmt"
	"github.com/ruslanDantsov/osmetrics-server/internal/model"
	"github.com/ruslanDantsov/osmetrics-server/internal/model/enum"
	"go.uber.org/zap"
	"sync"
)

type Storager interface {
	SaveMetric(m *model.Metrics) (*model.Metrics, error)
	GetMetric(metricID enum.MetricID) (*model.Metrics, bool)
	GetKnownMetrics() []string
}

type MemStorage struct {
	mu      sync.RWMutex
	Storage map[string]*model.Metrics
	Log     zap.Logger
}

func NewMemStorage(log zap.Logger) *MemStorage {
	return &MemStorage{
		Storage: make(map[string]*model.Metrics),
		Log:     log,
	}
}

func (s *MemStorage) GetKnownMetrics() []string {
	metricNames := make([]string, 0, len(s.Storage))
	for name := range s.Storage {
		metricNames = append(metricNames, name)
	}
	return metricNames
}

func (s *MemStorage) GetMetric(metricID enum.MetricID) (*model.Metrics, bool) {
	if val, found := s.Storage[metricID.String()]; found {
		if val.MType == "counter" {
			s.Log.Info(fmt.Sprintf("Get metric name=%v type=%v delta=%v", val.ID, val.MType, *val.Delta))
		}

		if val.MType == "gauge" {
			s.Log.Info(fmt.Sprintf("Get metric name=%v type=%v value=%v", val.ID, val.MType, *val.Value))
		}
		return val, true
	}

	return nil, false
}

func (s *MemStorage) SaveMetric(metric *model.Metrics) (*model.Metrics, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := metric.ID.String()
	existing, found := s.Storage[key]

	if metric.MType == "counter" {
		if found && existing.Delta != nil && metric.Delta != nil {
			*existing.Delta += *metric.Delta
			s.Log.Info(fmt.Sprintf("UPDATE counter_metric name=%v delta=%v", metric.ID, *existing.Delta))
			return existing, nil
		}
	}

	s.Storage[key] = metric

	if metric.MType == "counter" {
		s.Log.Info(fmt.Sprintf("SAVE %v metric id=%v delta=%v", metric.MType, metric.ID, *metric.Delta))
	}

	if metric.MType == "gauge" {
		s.Log.Info(fmt.Sprintf("SAVE %v metric id=%v value=%v", metric.MType, metric.ID, *metric.Value))
	}

	return metric, nil
}
