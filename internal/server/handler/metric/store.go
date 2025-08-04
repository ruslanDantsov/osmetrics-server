package metric

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/model"
	"github.com/ruslanDantsov/osmetrics-server/internal/server/constants"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

// MetricSaver предоставляет интерфейс для сохранения метрик.
type MetricSaver interface {
	SaveMetric(ctx context.Context, m *model.Metrics) (*model.Metrics, error)
	SaveAllMetrics(ctx context.Context, metricList model.MetricsList) (model.MetricsList, error)
}

// StoreMetricHandler обрабатывает HTTP-запросы на сохранение метрик.
type StoreMetricHandler struct {
	Storage MetricSaver
	Log     zap.Logger
}

// NewStoreMetricHandler создаёт новый экземпляр StoreMetricHandler.
func NewStoreMetricHandler(storage MetricSaver, log zap.Logger) *StoreMetricHandler {
	return &StoreMetricHandler{
		Storage: storage,
		Log:     log,
	}
}

// Store обрабатывает HTTP-запрос на сохранение одной метрики через URL-параметры.
//
// При успешной обработке возвращает HTTP 200 OK.
// В случае ошибок возвращает HTTP 400 с сообщением об ошибке.
func (h *StoreMetricHandler) Store(ginContext *gin.Context) {
	ctx := ginContext.Request.Context()

	metricType := strings.ToLower(ginContext.Param(constants.URLParamMetricType))
	metricName := ginContext.Param(constants.URLParamMetricName)
	metricValue := ginContext.Param(constants.URLParamMetricValue)

	var metricRequest *model.Metrics
	var err error

	switch metricType {
	case constants.CounterMetricType:
		metricRequest, err = model.NewMetricWithRawValues(constants.CounterMetricType, metricName, metricValue)
	case constants.GaugeMetricType:
		metricRequest, err = model.NewMetricWithRawValues(constants.GaugeMetricType, metricName, metricValue)
	default:
		h.Log.Warn(fmt.Sprintf("Metric type=%v is unsupported", metricType))
		metricRequest, err = model.NewMetricWithRawValues(metricType, metricName, metricValue)
	}

	if err != nil {
		h.Log.Error(err.Error())
		ginContext.String(http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.Storage.SaveMetric(ctx, metricRequest)
	if err != nil {
		h.Log.Error(err.Error())
		ginContext.String(http.StatusBadRequest, err.Error())
		return
	}

	ginContext.Status(http.StatusOK)
}
