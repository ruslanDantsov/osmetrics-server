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

	updatedMetric, err := h.Storage.SaveMetric(&metricRequest)
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
