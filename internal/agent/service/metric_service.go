package service

import (
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/ruslanDantsov/osmetrics-server/internal/agent/config"
	"github.com/ruslanDantsov/osmetrics-server/internal/agent/constants"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/model"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/model/enum"
	"github.com/ruslanDantsov/osmetrics-server/proto/metrics"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"math/rand"
	"runtime"
	"sync"
)

// MetricService отвечает за сбор и отправку метрик.
// Хранит текущее состояние метрик, использует клиент для HTTP-запросов и логгер.
type MetricService struct {
	mu      sync.Mutex
	log     *zap.Logger
	config  *config.AgentConfig
	metrics map[enum.MetricID]interface{}
	pubKey  *rsa.PublicKey
	localIP string

	conn   *grpc.ClientConn
	client metrics.MetricsServiceClient
}

// NewMetricService создает и возвращает новый экземпляр MetricService.
func NewMetricService(log *zap.Logger, agentConfig *config.AgentConfig, pubKey *rsa.PublicKey, localIP string) (*MetricService, error) {
	conn, err := grpc.Dial(agentConfig.Address, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	client := metrics.NewMetricsServiceClient(conn)

	return &MetricService{
		log:     log,
		config:  agentConfig,
		metrics: make(map[enum.MetricID]interface{}),
		pubKey:  pubKey,
		localIP: localIP,
		conn:    conn,
		client:  client,
	}, nil
}

func (ms *MetricService) Close() error {
	return ms.conn.Close()
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
	protoMetric := convertToProto(metric)
	_, err := ms.client.StoreMetric(ctx, protoMetric)
	if err != nil {
		return fmt.Errorf("failed to send metric via gRPC: %w", err)
	}
	return nil
}

func convertToProto(m model.Metrics) *metrics.Metric {
	metric := &metrics.Metric{
		Id:   string(m.ID),
		Type: m.MType,
	}
	if m.Delta != nil {
		metric.Delta = *m.Delta
	}
	if m.Value != nil {
		metric.Value = *m.Value
	}
	return metric
}
