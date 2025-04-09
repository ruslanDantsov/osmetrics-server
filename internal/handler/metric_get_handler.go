package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ruslanDantsov/osmetrics-server/internal/constants/server"
	"github.com/ruslanDantsov/osmetrics-server/internal/logging"
	"github.com/ruslanDantsov/osmetrics-server/internal/model/enum/metric"
	"github.com/ruslanDantsov/osmetrics-server/internal/repository"
	"net/http"
	"strconv"
)

type MetricGetHandler struct {
	Storage repository.Storager
	Log     logging.Logger
}

func NewMetricGetHandler(storage repository.Storager, log logging.Logger) *MetricGetHandler {
	return &MetricGetHandler{
		Storage: storage,
		Log:     log,
	}
}

func (h *MetricGetHandler) ServeHTTP(ginContext *gin.Context) {
	h.Log.Info(fmt.Sprintf("Handle request %v", ginContext.Request.RequestURI))

	metricType := ginContext.Param(server.URLParamMetricType)
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

func (h *MetricGetHandler) handleGetCounterMetric(ginContext *gin.Context) {
	rawMetricName := ginContext.Param(server.URLParamMetricName)

	metricName, err := metric.ParseMetricName(rawMetricName)
	if err != nil {
		h.Log.Error(err.Error())
		ginContext.String(http.StatusNotFound, "Metric name is unsupported")
		return
	}

	counterModel, found := h.Storage.GetCounterMetric(metricName)

	if !found {
		h.Log.Error(fmt.Sprintf("The counter_metric name=%v not found", metricName))
		ginContext.String(http.StatusNotFound, "Metric not found")
		return
	}

	ginContext.Header("Content-Type", "text/html")
	ginContext.String(http.StatusOK, strconv.FormatInt(counterModel.Value, 10))
}

func (h *MetricGetHandler) handleGetGaugeMetric(ginContext *gin.Context) {
	rawMetricName := ginContext.Param(server.URLParamMetricName)

	metricName, err := metric.ParseMetricName(rawMetricName)
	if err != nil {
		h.Log.Error(err.Error())
		ginContext.String(http.StatusNotFound, "Metric name is unsupported")
		return
	}

	gaugeModel, found := h.Storage.GetGaugeMetric(metricName)

	if !found {
		h.Log.Error(fmt.Sprintf("The gauge_metric name=%v not found", metricName))
		ginContext.String(http.StatusNotFound, "Metric not found")
		return
	}

	ginContext.Header("Content-Type", "text/html")
	ginContext.String(http.StatusOK, strconv.FormatFloat(gaugeModel.Value, 'f', -1, 64))
}
