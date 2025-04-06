package service

import (
	"fmt"
	"github.com/ruslanDantsov/osmetrics-server/internal/client"
	"github.com/ruslanDantsov/osmetrics-server/internal/logging"
	"github.com/ruslanDantsov/osmetrics-server/internal/model/enum/metric"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
)

type MetricService struct {
	mu      sync.Mutex
	Log     logging.Logger
	Client  client.RestClient
	Metrics map[metric.Metric]interface{}
}

func NewMetricService(log logging.Logger, client client.RestClient) *MetricService {
	return &MetricService{
		Log:     log,
		Client:  client,
		Metrics: make(map[metric.Metric]interface{}),
	}
}

func (ms *MetricService) CollectMetrics() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.appendMetric(metric.Alloc, float64(m.Alloc))
	ms.appendMetric(metric.BuckHashSys, float64(m.BuckHashSys))
	ms.appendMetric(metric.Frees, float64(m.Frees))
	ms.appendMetric(metric.GCCPUFraction, m.GCCPUFraction)
	ms.appendMetric(metric.GCSys, float64(m.GCSys))
	ms.appendMetric(metric.HeapAlloc, float64(m.HeapAlloc))
	ms.appendMetric(metric.HeapIdle, float64(m.HeapIdle))
	ms.appendMetric(metric.HeapInuse, float64(m.HeapInuse))
	ms.appendMetric(metric.HeapObjects, float64(m.HeapObjects))
	ms.appendMetric(metric.HeapReleased, float64(m.HeapReleased))
	ms.appendMetric(metric.HeapSys, float64(m.HeapSys))
	ms.appendMetric(metric.LastGC, float64(m.LastGC))
	ms.appendMetric(metric.Lookups, float64(m.Lookups))
	ms.appendMetric(metric.MCacheInuse, float64(m.MCacheInuse))
	ms.appendMetric(metric.MCacheSys, float64(m.MCacheSys))
	ms.appendMetric(metric.MSpanInuse, float64(m.MSpanInuse))
	ms.appendMetric(metric.MSpanSys, float64(m.MSpanSys))
	ms.appendMetric(metric.Mallocs, float64(m.Mallocs))
	ms.appendMetric(metric.NextGC, float64(m.NextGC))
	ms.appendMetric(metric.NumForcedGC, float64(m.NumForcedGC))
	ms.appendMetric(metric.NextGC, float64(m.NextGC))
	ms.appendMetric(metric.OtherSys, float64(m.OtherSys))
	ms.appendMetric(metric.PauseTotalNs, float64(m.PauseTotalNs))
	ms.appendMetric(metric.StackInuse, float64(m.StackInuse))
	ms.appendMetric(metric.StackSys, float64(m.StackSys))
	ms.appendMetric(metric.Sys, float64(m.Sys))
	ms.appendMetric(metric.TotalAlloc, float64(m.TotalAlloc))
	ms.aggregateMetric(metric.PollCount, 1)
	ms.appendMetric(metric.RandomValue, rand.Float64())
}

func (ms *MetricService) appendMetric(metricType metric.Metric, value float64) {
	ms.Metrics[metricType] = value
}

func (ms *MetricService) aggregateMetric(metricType metric.Metric, value int64) {
	if existingMetric, found := ms.Metrics[metricType]; found {
		ms.Metrics[metricType] = existingMetric.(int64) + value
	} else {
		ms.Metrics[metricType] = value
	}
}

func (ms *MetricService) SendMetrics() {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	for metricType, value := range ms.Metrics {
		switch metricType.Type {
		case metric.Gauge:
			err := ms.sendGaugeMetric(metricType.Name, value.(float64))
			if err != nil {
				ms.Log.Error(fmt.Sprintf("Failed to send metric %s: %v\n", metricType.Name, err))
				continue
			}
		case metric.Counter:
			err := ms.sendCounterMetric(metricType.Name, value.(int64))
			if err != nil {
				ms.Log.Error(fmt.Sprintf("Failed to send metric %s: %v\n", metricType.Name, err))
				continue
			}
			ms.Metrics[metricType] = int64(0)
		}
	}
}

func (ms *MetricService) sendGaugeMetric(name string, value float64) error {
	url := fmt.Sprintf("http://localhost:8080/update/gauge/%s/%f", name, value)
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
	url := fmt.Sprintf("http://localhost:8080/update/counter/%s/%v", name, value)
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
