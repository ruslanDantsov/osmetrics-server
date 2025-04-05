package repository

import (
	"fmt"
	"github.com/ruslanDantsov/osmetrics-server/internal/logging"
	"github.com/ruslanDantsov/osmetrics-server/internal/model"
	"github.com/ruslanDantsov/osmetrics-server/internal/model/enum/metric"
	"sync"
)

type Storager interface {
	SaveGaugeMetric(model *model.GaugeMetricModel) (*model.GaugeMetricModel, error)
	GetGaugeMetric(metricType metric.MetricType) (*model.GaugeMetricModel, bool)
	SaveCounterMetric(model *model.CounterMetricModel) (*model.CounterMetricModel, error)
	GetCounterMetric(metricType metric.MetricType) (*model.CounterMetricModel, bool)
}

type MemStorage struct {
	mu      sync.RWMutex
	Storage map[string]interface{}
	Log     logging.Logger
}

func NewMemStorage(log logging.Logger) *MemStorage {
	return &MemStorage{
		Storage: make(map[string]interface{}),
		Log:     log,
	}
}

func (s *MemStorage) GetGaugeMetric(metricType metric.MetricType) (*model.GaugeMetricModel, bool) {

	if valRaw, found := s.Storage[metricType.String()]; found {
		if val, ok := valRaw.(*model.GaugeMetricModel); ok {
			s.Log.Info(fmt.Sprintf("GET gauge_metric type=%v value=%v", val.MetricType, val.Value))
			return val, true
		}
	}

	return nil, false
}

func (s *MemStorage) GetCounterMetric(metricType metric.MetricType) (*model.CounterMetricModel, bool) {

	if valRaw, found := s.Storage[metricType.String()]; found {
		if val, ok := valRaw.(*model.CounterMetricModel); ok {
			s.Log.Info(fmt.Sprintf("GET counter_metric type=%v value=%v", val.MetricType, val.Value))
			return val, true
		}
	}

	return nil, false
}

func (s *MemStorage) SaveGaugeMetric(model *model.GaugeMetricModel) (*model.GaugeMetricModel, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Storage[model.MetricType.String()] = model
	s.Log.Info(fmt.Sprintf("SAVE gauge_metric type=%v value=%v", model.MetricType, model.Value))

	return model, nil
}

func (s *MemStorage) SaveCounterMetric(rawModel *model.CounterMetricModel) (*model.CounterMetricModel, error) {
	existingCounterModel, found := s.GetCounterMetric(rawModel.MetricType)

	s.mu.Lock()
	defer s.mu.Unlock()

	if found {
		existingCounterModel.Value += rawModel.Value
		s.Log.Info(fmt.Sprintf("UPDATE counter_metric type=%v value=%v", existingCounterModel.MetricType, existingCounterModel.Value))
	} else {
		s.Storage[rawModel.MetricType.String()] = rawModel
		s.Log.Info(fmt.Sprintf("SAVE counter_metric type=%v value=%v", rawModel.MetricType, rawModel.Value))
	}

	updatedModel := s.Storage[rawModel.MetricType.String()]
	return updatedModel.(*model.CounterMetricModel), nil
}
