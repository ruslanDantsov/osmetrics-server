package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ruslanDantsov/osmetrics-server/internal/logging"
	"github.com/ruslanDantsov/osmetrics-server/internal/model"
	"github.com/ruslanDantsov/osmetrics-server/internal/repository"
	"net/http"
)

type MetricPostHandler struct {
	Storage repository.Storager
	Log     logging.Logger
}

func NewMetricPostHandler(storage repository.Storager, log logging.Logger) *MetricPostHandler {
	return &MetricPostHandler{
		Storage: storage,
		Log:     log,
	}
}

func (h *MetricPostHandler) ServeHTTP(ginContext *gin.Context) {
	h.Log.Info(fmt.Sprintf("Handle request %v", ginContext.Request.RequestURI))

	metricType := ginContext.Param("type")
	switch metricType {
	case "gauge":
		h.handlePostGaugeMetric(ginContext)
	case "counter":
		h.handlePostCounterMetric(ginContext)
	default:
		h.Log.Error(fmt.Sprintf("Metric type=%v is unsupported", metricType))
		ginContext.String(http.StatusBadRequest, "Metric type is unsupported")
	}

}

func (h *MetricPostHandler) handlePostCounterMetric(ginContext *gin.Context) {
	metricName := ginContext.Param("name")
	metricValue := ginContext.Param("value")

	counterModel, err := model.NewCounterMetricModelWithRawValues(metricName, metricValue)
	if err != nil {
		h.Log.Error(err.Error())
		ginContext.String(http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.Storage.SaveCounterMetric(counterModel)
	if err != nil {
		h.Log.Error(err.Error())
		ginContext.String(http.StatusBadRequest, err.Error())
		return
	}

	ginContext.Status(http.StatusOK)
}

func (h *MetricPostHandler) handlePostGaugeMetric(ginContext *gin.Context) {
	metricName := ginContext.Param("name")
	metricValue := ginContext.Param("value")

	gaugeModel, err := model.NewGaugeMetricModelWithRawValues(metricName, metricValue)
	if err != nil {
		h.Log.Error(err.Error())
		ginContext.String(http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.Storage.SaveGaugeMetric(gaugeModel)
	if err != nil {
		h.Log.Error(err.Error())
		ginContext.String(http.StatusBadRequest, err.Error())
		return
	}

	ginContext.Status(http.StatusOK)
}
