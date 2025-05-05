package metric

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"
	"github.com/ruslanDantsov/osmetrics-server/internal/constants"
	"github.com/ruslanDantsov/osmetrics-server/internal/model"
	"net/http"
	"strings"
)

func (h *MetricHandler) GetJSON(ginContext *gin.Context) {
	var metricRequest model.Metrics

	if err := easyjson.UnmarshalFromReader(ginContext.Request.Body, &metricRequest); err != nil {
		h.Log.Error(fmt.Sprintf("Error on unmarshal data from request. %v", err))
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if strings.ToLower(metricRequest.MType) != constants.GaugeMetricType && strings.ToLower(metricRequest.MType) != constants.CounterMetricType {
		h.Log.Warn(fmt.Sprintf("Metric type=%v is unsupported", metricRequest.MType))
	}

	existingMetric, found := h.Storage.GetMetric(metricRequest.ID)
	if !found {
		h.Log.Warn(fmt.Sprintf("The metric ID=%v not found", metricRequest.ID))
		ginContext.JSON(http.StatusNotFound, gin.H{"error": "Metric not found"})
		return
	}

	ginContext.Header("Content-Type", "application/json")
	ginContext.Writer.WriteHeader(http.StatusOK)

	_, err := easyjson.MarshalToWriter(existingMetric, ginContext.Writer)

	if existingMetric.MType == constants.CounterMetricType {
		h.Log.Debug(fmt.Sprintf("Return: Metric ID=%v Value=%v", existingMetric.ID, existingMetric.Delta))
	}

	if existingMetric.MType == constants.CounterMetricType {
		h.Log.Debug(fmt.Sprintf("Return: Metric ID=%v Value=%v", existingMetric.ID, existingMetric.Value))
	}

	if err != nil {
		h.Log.Error(fmt.Sprintf("Error on marshal metric data. %v", err))
		ginContext.JSON(http.StatusNotFound, gin.H{"error": "Metric not found"})
	}
}
