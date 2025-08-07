// Package metric provides HTTP handlers for system metrics.
package metric

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/model"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/model/enum"
	"github.com/ruslanDantsov/osmetrics-server/internal/server/constants"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// MetricGetter определяет интерфейс для получения метрик по их идентификатору.
type MetricGetter interface {
	GetMetric(ctx context.Context, metricID enum.MetricID) (*model.Metrics, bool)
	GetKnownMetrics(ctx context.Context) []string
}

// GetMetricHandler представляет обработчик HTTP-запросов для получения метрик.
type GetMetricHandler struct {
	Storage MetricGetter
	Log     zap.Logger
}

// NewGetMetricHandler создаёт новый экземпляр GetMetricHandler.
func NewGetMetricHandler(storage MetricGetter, log zap.Logger) *GetMetricHandler {
	return &GetMetricHandler{
		Storage: storage,
		Log:     log,
	}
}

// Get обрабатывает входящий HTTP-запрос для получения значения метрики.
//
// Определяет тип метрики (gauge или counter) и вызывает соответствующий обработчик.
// В случае некорректного типа возвращает ошибку 400.
func (h *GetMetricHandler) Get(ginContext *gin.Context) {
	metricType := ginContext.Param(constants.URLParamMetricType)
	switch metricType {
	case constants.GaugeMetricType:
		h.handleGetGaugeMetric(ginContext)
	case constants.CounterMetricType:
		h.handleGetCounterMetric(ginContext)
	default:
		h.Log.Error(fmt.Sprintf("Metric type=%v is unsupported", metricType))
		ginContext.String(http.StatusBadRequest, "Metric type is unsupported")
	}
}

func (h *GetMetricHandler) handleGetCounterMetric(ginContext *gin.Context) {
	ctx := ginContext.Request.Context()

	rawMetricID := ginContext.Param(constants.URLParamMetricName)

	metricID, err := enum.ParseMetricID(rawMetricID)
	if err != nil {
		h.Log.Error(err.Error())
		ginContext.String(http.StatusNotFound, "Metric name is unsupported")
		return
	}

	metricModel, found := h.Storage.GetMetric(ctx, metricID)

	if !found {
		h.Log.Warn(fmt.Sprintf("The counter_metric name=%v not found", metricID))
		ginContext.String(http.StatusNotFound, "Metric not found")
		return
	}

	ginContext.Header("Content-Type", "text/html")
	ginContext.String(http.StatusOK, strconv.FormatInt(*metricModel.Delta, 10))
}

func (h *GetMetricHandler) handleGetGaugeMetric(ginContext *gin.Context) {
	ctx := ginContext.Request.Context()

	rawMetricID := ginContext.Param(constants.URLParamMetricName)

	metricID, err := enum.ParseMetricID(rawMetricID)
	if err != nil {
		h.Log.Error(err.Error())
		ginContext.String(http.StatusNotFound, "Metric name is unsupported")
		return
	}

	gaugeModel, found := h.Storage.GetMetric(ctx, metricID)

	if !found {
		h.Log.Warn(fmt.Sprintf("The gauge_metric name=%v not found", metricID))
		ginContext.String(http.StatusNotFound, "Metric not found")
		return
	}

	ginContext.Header("Content-Type", "text/html")
	ginContext.String(http.StatusOK, strconv.FormatFloat(*gaugeModel.Value, 'f', -1, 64))
}
