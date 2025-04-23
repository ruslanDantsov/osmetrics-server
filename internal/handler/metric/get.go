package metric

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ruslanDantsov/osmetrics-server/internal/constants"
	"github.com/ruslanDantsov/osmetrics-server/internal/model/enum"
	"github.com/ruslanDantsov/osmetrics-server/internal/repository"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type MetricHandler struct {
	Storage repository.Storager
	Log     zap.Logger
}

func NewMetricHandler(storage repository.Storager, log zap.Logger) *MetricHandler {
	return &MetricHandler{
		Storage: storage,
		Log:     log,
	}
}

func (h *MetricHandler) Get(ginContext *gin.Context) {
	metricType := ginContext.Param(constants.URLParamMetricType)
	switch metricType {
	case "gauge":
		h.handleGetGaugeMetric(ginContext)
	case "counter":
		h.handleGetCounterMetric(ginContext)
	default:
		h.Log.Error(fmt.Sprintf("Metric type=%v is unsupported", metricType))
		ginContext.String(http.StatusBadRequest, "Metric type is unsupported")
	}
}

func (h *MetricHandler) handleGetCounterMetric(ginContext *gin.Context) {
	rawMetricName := ginContext.Param(constants.URLParamMetricName)

	metricName, err := enum.ParseMetricId(rawMetricName)
	if err != nil {
		h.Log.Error(err.Error())
		ginContext.String(http.StatusNotFound, "Metric name is unsupported")
		return
	}

	metricModel, found := h.Storage.GetMetric(metricName)

	if !found {
		h.Log.Error(fmt.Sprintf("The counter_metric name=%v not found", metricName))
		ginContext.String(http.StatusNotFound, "Metric not found")
		return
	}

	ginContext.Header("Content-Type", "text/html")
	ginContext.String(http.StatusOK, strconv.FormatInt(*metricModel.Delta, 10))
}

func (h *MetricHandler) handleGetGaugeMetric(ginContext *gin.Context) {
	rawMetricName := ginContext.Param(constants.URLParamMetricName)

	metricName, err := enum.ParseMetricId(rawMetricName)
	if err != nil {
		h.Log.Error(err.Error())
		ginContext.String(http.StatusNotFound, "Metric name is unsupported")
		return
	}

	gaugeModel, found := h.Storage.GetMetric(metricName)

	if !found {
		h.Log.Error(fmt.Sprintf("The gauge_metric name=%v not found", metricName))
		ginContext.String(http.StatusNotFound, "Metric not found")
		return
	}

	ginContext.Header("Content-Type", "text/html")
	ginContext.String(http.StatusOK, strconv.FormatFloat(*gaugeModel.Value, 'f', -1, 64))
}
