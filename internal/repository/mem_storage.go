package repository

import (
	"fmt"
	"github.com/ruslanDantsov/osmetrics-server/internal/logging"
	"github.com/ruslanDantsov/osmetrics-server/internal/model"
	"github.com/ruslanDantsov/osmetrics-server/internal/model/enum/metric"
	"sync"
)

//TODO: split logic for repository and service layers

type Storager interface {
	SaveGaugeMetric(model *model.GaugeMetricModel) (*model.GaugeMetricModel, error)
	GetGaugeMetric(metricType metric.Metric) (*model.GaugeMetricModel, bool)
	SaveCounterMetric(model *model.CounterMetricModel) (*model.CounterMetricModel, error)
	GetCounterMetric(metricType metric.Metric) (*model.CounterMetricModel, bool)
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

func (s *MemStorage) GetGaugeMetric(metricName metric.Metric) (*model.GaugeMetricModel, bool) {

	if valRaw, found := s.Storage[metricName.String()]; found {
		if val, ok := valRaw.(*model.GaugeMetricModel); ok {
			s.Log.Info(fmt.Sprintf("GET gauge_metric name=%v value=%v", val.Name, val.Value))
			return val, true
		}
	}

	return nil, false
}

func (s *MemStorage) GetCounterMetric(metricName metric.Metric) (*model.CounterMetricModel, bool) {

	if valRaw, found := s.Storage[metricName.String()]; found {
		if val, ok := valRaw.(*model.CounterMetricModel); ok {
			s.Log.Info(fmt.Sprintf("GET counter_metric name=%v value=%v", val.Name, val.Value))
			return val, true
		}
	}

	return nil, false
}

func (s *MemStorage) SaveGaugeMetric(model *model.GaugeMetricModel) (*model.GaugeMetricModel, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Storage[model.Name.String()] = model
	s.Log.Info(fmt.Sprintf("SAVE gauge_metric name=%v value=%v", model.Name, model.Value))

	return model, nil
}

func (s *MemStorage) SaveCounterMetric(rawModel *model.CounterMetricModel) (*model.CounterMetricModel, error) {
	existingCounterModel, found := s.GetCounterMetric(rawModel.Name)

	s.mu.Lock()
	defer s.mu.Unlock()

	if found {
		existingCounterModel.Value += rawModel.Value
		s.Log.Info(fmt.Sprintf("UPDATE counter_metric name=%v value=%v", existingCounterModel.Name, existingCounterModel.Value))
	} else {
		s.Storage[rawModel.Name.String()] = rawModel
		s.Log.Info(fmt.Sprintf("SAVE counter_metric name=%v value=%v", rawModel.Name, rawModel.Value))
	}

	updatedModel := s.Storage[rawModel.Name.String()]
	return updatedModel.(*model.CounterMetricModel), nil
}
