package metric

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/model"
	"github.com/ruslanDantsov/osmetrics-server/internal/server/constants"
	"net/http"
	"strconv"
	"strings"
)

func (h *StoreMetricHandler) StoreJSON(ginContext *gin.Context) {
	var metricRequest model.Metrics

	if err := easyjson.UnmarshalFromReader(ginContext.Request.Body, &metricRequest); err != nil {
		h.Log.Error(fmt.Sprintf("Error on unmarshal data from request. %v", err))
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if strings.ToLower(metricRequest.MType) != constants.GaugeMetricType && strings.ToLower(metricRequest.MType) != constants.CounterMetricType {
		h.Log.Warn(fmt.Sprintf("Metric type=%v is unsupported", metricRequest.MType))
	}

	ctx := ginContext.Request.Context()
	updatedMetric, err := h.Storage.SaveMetric(ctx, &metricRequest)
	if err != nil {
		h.Log.Error(err.Error())
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Can't update metric", "description": err.Error()})
		return
	}

	if _, err = easyjson.MarshalToWriter(updatedMetric, ginContext.Writer); err != nil {
		h.Log.Error(fmt.Sprintf("Error on marshal metric data. %v", err))
		ginContext.JSON(http.StatusNotFound, gin.H{"error": "Can't convert data to JSON"})
		return
	}

	ginContext.Status(http.StatusOK)
}

func (h *StoreMetricHandler) StoreBatchJSON(ginContext *gin.Context) {
	ctx := ginContext.Request.Context()
	var metrics []model.Metrics
	if err := ginContext.ShouldBindJSON(&metrics); err != nil {
		h.Log.Error("Failed to parse metrics batch: " + err.Error())
		ginContext.String(http.StatusBadRequest, "invalid request body")
		return
	}

	if len(metrics) == 0 {
		ginContext.Status(http.StatusOK)
		return
	}

	var metricsList []model.Metrics
	for _, metric := range metrics {
		var value string

		switch metric.MType {
		case constants.CounterMetricType:
			value = strconv.FormatInt(*metric.Delta, 10)
		case constants.GaugeMetricType:
			value = strconv.FormatFloat(*metric.Value, 'f', -1, 64)
		default:
			continue
		}

		rawMmetric, err := model.NewMetricWithRawValues(metric.MType, string(metric.ID), value)
		if err != nil {
			ginContext.String(http.StatusBadRequest, "Error on model construction")
			return
		}

		metricsList = append(metricsList, *rawMmetric)
	}

	_, err := h.Storage.SaveAllMetrics(ctx, metricsList)
	if err != nil {
		ginContext.String(http.StatusBadRequest, "Error on saving batch metrics")
		return
	}

	ginContext.Status(http.StatusOK)

}
