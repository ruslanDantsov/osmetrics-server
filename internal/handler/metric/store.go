package metric

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ruslanDantsov/osmetrics-server/internal/constants"
	"github.com/ruslanDantsov/osmetrics-server/internal/model"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type MetricSaver interface {
	SaveMetric(ctx context.Context, m *model.Metrics) (*model.Metrics, error)
}

type StoreMetricHandler struct {
	Storage MetricSaver
	Log     zap.Logger
}

func NewStoreMetricHandler(storage MetricSaver, log zap.Logger) *StoreMetricHandler {
	return &StoreMetricHandler{
		Storage: storage,
		Log:     log,
	}
}

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
