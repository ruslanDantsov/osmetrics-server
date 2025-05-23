package service

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/ruslanDantsov/osmetrics-server/internal/config"
	"github.com/ruslanDantsov/osmetrics-server/internal/constants"
	"github.com/ruslanDantsov/osmetrics-server/internal/model"
	"github.com/ruslanDantsov/osmetrics-server/internal/model/enum"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
)

const (
	UpdateMetricURL  = "http://%v/update"
	UpdateMetricsURL = "http://%v/updates"
)

type RestClient interface {
	R() *resty.Request
}

type MetricService struct {
	mu      sync.Mutex
	Log     *zap.Logger
	Client  RestClient
	config  *config.AgentConfig
	Metrics map[enum.MetricID]interface{}
}

func NewMetricService(log *zap.Logger, client RestClient, agentConfig *config.AgentConfig) *MetricService {
	return &MetricService{
		Log:     log,
		Client:  client,
		config:  agentConfig,
		Metrics: make(map[enum.MetricID]interface{}),
	}
}
func (ms *MetricService) CollectMetrics() {
	ms.Log.Info("Collecting metrics...")
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	ms.mu.Lock()
	defer ms.mu.Unlock()

	// Define metrics to collect as map
	metrics := map[enum.MetricID]float64{
		enum.Alloc:         float64(memStats.Alloc),
		enum.BuckHashSys:   float64(memStats.BuckHashSys),
		enum.Frees:         float64(memStats.Frees),
		enum.GCCPUFraction: memStats.GCCPUFraction,
		enum.GCSys:         float64(memStats.GCSys),
		enum.HeapAlloc:     float64(memStats.HeapAlloc),
		enum.HeapIdle:      float64(memStats.HeapIdle),
		enum.HeapInuse:     float64(memStats.HeapInuse),
		enum.HeapObjects:   float64(memStats.HeapObjects),
		enum.HeapReleased:  float64(memStats.HeapReleased),
		enum.HeapSys:       float64(memStats.HeapSys),
		enum.LastGC:        float64(memStats.LastGC),
		enum.Lookups:       float64(memStats.Lookups),
		enum.MCacheInuse:   float64(memStats.MCacheInuse),
		enum.MCacheSys:     float64(memStats.MCacheSys),
		enum.MSpanInuse:    float64(memStats.MSpanInuse),
		enum.MSpanSys:      float64(memStats.MSpanSys),
		enum.Mallocs:       float64(memStats.Mallocs),
		enum.NextGC:        float64(memStats.NextGC),
		enum.OtherSys:      float64(memStats.OtherSys),
		enum.PauseTotalNs:  float64(memStats.PauseTotalNs),
		enum.StackInuse:    float64(memStats.StackInuse),
		enum.StackSys:      float64(memStats.StackSys),
		enum.Sys:           float64(memStats.Sys),
		enum.TotalAlloc:    float64(memStats.TotalAlloc),
		enum.NumForcedGC:   float64(memStats.NumForcedGC),
		enum.NumGC:         float64(memStats.NumGC),
		enum.RandomValue:   rand.Float64(),
	}

	for id, value := range metrics {
		ms.appendMetric(id, value)
	}

	ms.aggregateMetric(enum.PollCount, 1)
}

func (ms *MetricService) appendMetric(metricType enum.MetricID, value float64) {
	ms.Metrics[metricType] = value
}

func (ms *MetricService) aggregateMetric(metricType enum.MetricID, value int64) {
	if existingMetric, found := ms.Metrics[metricType]; found {
		ms.Metrics[metricType] = existingMetric.(int64) + value
	} else {
		ms.Metrics[metricType] = value
	}
}

func (ms *MetricService) sendMetric(ID enum.MetricID, mType string, value interface{}) error {
	url := fmt.Sprintf(UpdateMetricURL, ms.config.Address)

	metric := model.Metrics{
		ID:    ID,
		MType: mType,
	}

	switch mType {
	case constants.GaugeMetricType:
		v := value.(float64)
		metric.Value = &v
	case constants.CounterMetricType:
		v := value.(int64)
		metric.Delta = &v
	}

	json, err := metric.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal metric: %w", err)
	}

	resp, err := ms.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(json).
		Post(url)

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("bad response for metric %s: %v", ID, resp.StatusCode())
	}

	return nil
}

func (ms *MetricService) SendMetrics() {
	ms.Log.Info("Sending metrics...")

	for metricID, genericValue := range ms.Metrics {
		var err error
		switch value := genericValue.(type) {
		case float64:
			err = ms.sendMetric(metricID, constants.GaugeMetricType, value)
		case int64:
			err = ms.sendMetric(metricID, constants.CounterMetricType, value)
			if err == nil {
				ms.Metrics[metricID] = int64(0)
			}
		}
		if err != nil {
			ms.Log.Error(fmt.Sprintf("Failed to send metric %s: %v\n", metricID, err))
		}
	}
}

func (ms *MetricService) SendAllMetrics() {
	ms.Log.Info("Sending batch of metrics...")
	url := fmt.Sprintf(UpdateMetricsURL, ms.config.Address)
	var metricList []model.Metrics

	for metricID, genericValue := range ms.Metrics {
		metric := model.Metrics{
			ID: metricID,
		}

		switch value := genericValue.(type) {
		case float64:
			metric.MType = constants.GaugeMetricType
			metric.Value = &value
		case int64:
			metric.MType = constants.CounterMetricType
			metric.Delta = &value
			ms.Metrics[metricID] = int64(0)
		}
		metricList = append(metricList, metric)
	}

	json, err := model.MetricsList(metricList).MarshalJSON()
	if err != nil {
		ms.Log.Error("failed to marshal batch of metrics: %w", zap.Error(err))
	}

	resp, err := ms.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(json).
		Post(url)

	if err != nil {
		ms.Log.Error("failed to send batch of metrics: %w", zap.Error(err))
	}

	if resp.StatusCode() != http.StatusOK {
		ms.Log.Error(fmt.Sprintf("bad response for batch of metrics %v", resp.StatusCode()))
	}
}
