package service

import (
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/ruslanDantsov/osmetrics-server/internal/agent/config"
	"github.com/ruslanDantsov/osmetrics-server/internal/agent/constants"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/crypto"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/model"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/model/enum"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
)

// RestClient определяет интерфейс для HTTP-клиента, совместимого с resty.
type RestClient interface {
	R() *resty.Request
}

// MetricService отвечает за сбор и отправку метрик.
// Хранит текущее состояние метрик, использует клиент для HTTP-запросов и логгер.
type MetricService struct {
	mu      sync.Mutex
	log     *zap.Logger
	client  RestClient
	config  *config.AgentConfig
	metrics map[enum.MetricID]interface{}
	pubKey  *rsa.PublicKey
}

// NewMetricService создает и возвращает новый экземпляр MetricService.
func NewMetricService(log *zap.Logger, client RestClient, agentConfig *config.AgentConfig, pubKey *rsa.PublicKey) *MetricService {
	return &MetricService{
		log:     log,
		client:  client,
		config:  agentConfig,
		metrics: make(map[enum.MetricID]interface{}),
		pubKey:  pubKey,
	}
}

// CollectMetrics собирает метрики из runtime и отправляет их в канал metricChan.
func (ms *MetricService) CollectMetrics(metricChan chan<- model.Metrics) {
	ms.log.Info("Collecting metrics...")
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	ms.mu.Lock()
	defer ms.mu.Unlock()

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
		metricChan <- model.Metrics{
			ID:    id,
			MType: constants.GaugeMetricType,
			Value: &value,
		}
	}

	count := int64(1)
	metricChan <- model.Metrics{
		ID:    enum.PollCount,
		MType: constants.CounterMetricType,
		Delta: &count,
	}
}

// CollectAdditionalMetrics собирает дополнительные системные метрики,
// включая свободную и общую память, а также количество CPU, и отправляет их в канал.
func (ms *MetricService) CollectAdditionalMetrics(metricChan chan<- model.Metrics) {
	ms.log.Info("Collecting additional metrics...")
	memInfo, _ := mem.VirtualMemory()
	cpuCount, _ := cpu.Counts(false)

	ms.mu.Lock()
	defer ms.mu.Unlock()

	metrics := map[enum.MetricID]float64{
		enum.FreeMemory:      float64(memInfo.Free),
		enum.TotalMemory:     float64(memInfo.Total),
		enum.CPUutilization1: float64(cpuCount),
	}

	for id, value := range metrics {
		metricChan <- model.Metrics{
			ID:    id,
			MType: constants.GaugeMetricType,
			Value: &value,
		}
	}
}

// Worker запускает воркер, который читает метрики из канала metricChan
// и отправляет их на сервер. Завершается при отмене контекста.
func (ms *MetricService) Worker(ctx context.Context, metricChan chan model.Metrics) {
	for {
		select {
		case <-ctx.Done():
			ms.log.Info("Worker received shutdown signal")
			return
		case metric, ok := <-metricChan:
			if !ok {
				ms.log.Info("Metric channel closed, worker exiting")
				return
			}

			if err := ms.sendMetric(ctx, metric); err != nil {
				ms.log.Error("Failed to send metric", zap.String("id", string(metric.ID)), zap.Error(err))
			}
		}
	}
}

func (ms *MetricService) sendMetric(ctx context.Context, metric model.Metrics) error {
	url := fmt.Sprintf(constants.UpdateMetricURL, ms.config.Address)

	json, err := metric.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal metric: %w", err)
	}

	payload, err := crypto.EncryptPayload(ms.pubKey, json)
	if err != nil {
		return fmt.Errorf("failed to encrypt metric: %w", err)
	}

	resp, err := ms.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		SetContext(ctx).
		Post(url)

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("bad response for SendingMetric %s: %v", metric.ID, resp.StatusCode())
	}

	return nil
}

// SendAllMetrics отправляет все накопленные метрики в виде батча на сервер.
func (ms *MetricService) SendAllMetrics() {
	ms.log.Info("Sending batch of metrics...")
	url := fmt.Sprintf(constants.UpdateMetricsURL, ms.config.Address)
	var metricList []model.Metrics

	ms.mu.Lock()
	for metricID, genericValue := range ms.metrics {
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
			ms.metrics[metricID] = int64(0)
		}
		metricList = append(metricList, metric)
	}
	ms.mu.Unlock()

	json, err := model.MetricsList(metricList).MarshalJSON()
	if err != nil {
		ms.log.Error("failed to marshal batch of metrics: %w", zap.Error(err))
	}

	resp, err := ms.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(json).
		Post(url)

	if err != nil {
		ms.log.Error("failed to send batch of metrics: %w", zap.Error(err))
	}

	if resp.StatusCode() != http.StatusOK {
		ms.log.Error(fmt.Sprintf("bad response for batch of metrics %v", resp.StatusCode()))
	}
}
