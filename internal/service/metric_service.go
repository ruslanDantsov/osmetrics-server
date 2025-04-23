package service

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/ruslanDantsov/osmetrics-server/internal/config"
	"github.com/ruslanDantsov/osmetrics-server/internal/logger"
	"github.com/ruslanDantsov/osmetrics-server/internal/model/enum"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
)

type RestClient interface {
	R() *resty.Request
}

type MetricService struct {
	mu      sync.Mutex
	Log     logger.Logger
	Client  RestClient
	config  *config.AgentConfig
	Metrics map[enum.MetricID]interface{}
}

func NewMetricService(log logger.Logger, client RestClient, agentConfig *config.AgentConfig) *MetricService {
	return &MetricService{
		Log:     log,
		Client:  client,
		config:  agentConfig,
		Metrics: make(map[enum.MetricID]interface{}),
	}
}

func (ms *MetricService) CollectMetrics() {
	ms.Log.Info("Start process for collecting metrics")
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.appendMetric(enum.Alloc, float64(memStats.Alloc))
	ms.appendMetric(enum.BuckHashSys, float64(memStats.BuckHashSys))
	ms.appendMetric(enum.Frees, float64(memStats.Frees))
	ms.appendMetric(enum.GCCPUFraction, memStats.GCCPUFraction)
	ms.appendMetric(enum.GCSys, float64(memStats.GCSys))
	ms.appendMetric(enum.HeapAlloc, float64(memStats.HeapAlloc))
	ms.appendMetric(enum.HeapIdle, float64(memStats.HeapIdle))
	ms.appendMetric(enum.HeapInuse, float64(memStats.HeapInuse))
	ms.appendMetric(enum.HeapObjects, float64(memStats.HeapObjects))
	ms.appendMetric(enum.HeapReleased, float64(memStats.HeapReleased))
	ms.appendMetric(enum.HeapSys, float64(memStats.HeapSys))
	ms.appendMetric(enum.LastGC, float64(memStats.LastGC))
	ms.appendMetric(enum.Lookups, float64(memStats.Lookups))
	ms.appendMetric(enum.MCacheInuse, float64(memStats.MCacheInuse))
	ms.appendMetric(enum.MCacheSys, float64(memStats.MCacheSys))
	ms.appendMetric(enum.MSpanInuse, float64(memStats.MSpanInuse))
	ms.appendMetric(enum.MSpanSys, float64(memStats.MSpanSys))
	ms.appendMetric(enum.Mallocs, float64(memStats.Mallocs))
	ms.appendMetric(enum.NextGC, float64(memStats.NextGC))
	ms.appendMetric(enum.NumForcedGC, float64(memStats.NumForcedGC))
	ms.appendMetric(enum.NextGC, float64(memStats.NextGC))
	ms.appendMetric(enum.OtherSys, float64(memStats.OtherSys))
	ms.appendMetric(enum.PauseTotalNs, float64(memStats.PauseTotalNs))
	ms.appendMetric(enum.StackInuse, float64(memStats.StackInuse))
	ms.appendMetric(enum.StackSys, float64(memStats.StackSys))
	ms.appendMetric(enum.Sys, float64(memStats.Sys))
	ms.appendMetric(enum.TotalAlloc, float64(memStats.TotalAlloc))
	ms.aggregateMetric(enum.PollCount, 1)
	ms.appendMetric(enum.RandomValue, rand.Float64())
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

func (ms *MetricService) SendMetrics() {
	ms.Log.Info("Start process for sending metrics")

	for metricID, genericValue := range ms.Metrics {
		switch value := genericValue.(type) {
		case float64:
			err := ms.sendGaugeMetric(metricID.String(), value)
			if err != nil {
				ms.Log.Error(fmt.Sprintf("Failed to send metric %s: %v\n", metricID, err))
				continue
			}
		case int64:
			err := ms.sendCounterMetric(metricID.String(), value)
			if err != nil {
				ms.Log.Error(fmt.Sprintf("Failed to send metric %s: %v\n", metricID, err))
				continue
			}
			ms.Metrics[metricID] = int64(0)
		}
	}
}

func (ms *MetricService) sendGaugeMetric(name string, value float64) error {
	url := fmt.Sprintf("http://%v/update/gauge/%s/%f", ms.config.Address, name, value)
	resp, err := ms.Client.R().
		SetHeader("Content-Type", "text/plain").
		Post(url)

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("bad response for metric %s: %v", name, resp.StatusCode())
	}

	return nil
}

func (ms *MetricService) sendCounterMetric(name string, value int64) error {
	url := fmt.Sprintf("http://%v/update/counter/%s/%v", ms.config.Address, name, value)
	resp, err := ms.Client.R().
		SetHeader("Content-Type", "text/plain").
		Post(url)

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("bad response for metric %s: %v", name, resp.StatusCode())
	}

	return nil
}
