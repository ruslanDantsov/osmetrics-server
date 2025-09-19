package handler

import (
	"context"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/model"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/model/enum"
	"github.com/ruslanDantsov/osmetrics-server/proto/metrics"
	"go.uber.org/zap"
)

type MetricSaver interface {
	SaveMetric(ctx context.Context, m *model.Metrics) (*model.Metrics, error)
	SaveAllMetrics(ctx context.Context, metricList model.MetricsList) (model.MetricsList, error)
}

// StoreMetricHandler обрабатывает HTTP-запросы на сохранение метрик.
type StoreMetricHandler struct {
	Storage MetricSaver
	Log     zap.Logger
}
type ServerMetricsHandler struct {
	metrics.UnimplementedMetricsServiceServer
	storage MetricSaver
	logger  *zap.Logger
}

// NewServerMetricsHandler создаёт новый хэндлер
func NewServerMetricsHandler(storage MetricSaver, logger *zap.Logger) *ServerMetricsHandler {
	return &ServerMetricsHandler{
		storage: storage,
		logger:  logger,
	}
}

// StoreMetric сохраняет одну метрику
func (h *ServerMetricsHandler) StoreMetric(ctx context.Context, req *metrics.Metric) (*metrics.MetricResponse, error) {
	if req == nil {
		return &metrics.MetricResponse{Success: false, Error: "metric is nil"}, nil
	}
	metricId, err := enum.ParseMetricID(req.Id)
	if err != nil {
		h.logger.Error("invalid metric ID", zap.String("id", req.Id), zap.Error(err))
		return &metrics.MetricResponse{Success: false, Error: "invalid metric ID"}, nil
	}
	var metricRequest = model.Metrics{
		ID:    metricId,
		MType: req.Type,
		Delta: &req.Delta,
		Value: &req.Value,
	}

	_, err = h.storage.SaveMetric(ctx, &metricRequest)
	if err != nil {
		h.logger.Error("failed to store metric", zap.Error(err))
		return &metrics.MetricResponse{Success: false, Error: err.Error()}, nil
	}

	return &metrics.MetricResponse{Success: true}, nil
}

// StoreMetricsBatch сохраняет пачку метрик
func (h *ServerMetricsHandler) StoreMetricsBatch(ctx context.Context, req *metrics.MetricsBatchRequest) (*metrics.MetricsBatchResponse, error) {
	if req == nil || len(req.Metrics) == 0 {
		return &metrics.MetricsBatchResponse{Success: true}, nil // ничего сохранять не нужно
	}

	return &metrics.MetricsBatchResponse{Success: true}, nil
}
