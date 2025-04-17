package metric

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ruslanDantsov/osmetrics-server/internal/constants"
	"github.com/ruslanDantsov/osmetrics-server/internal/model"
	"net/http"
)

func (h *MetricHandler) Create(ginContext *gin.Context) {
	metricType := ginContext.Param(constants.URLParamMetricType)
	switch metricType {
	case "gauge":
		h.handleCreateGaugeMetric(ginContext)
	case "counter":
		h.handleCreateCounterMetric(ginContext)
	default:
		h.Log.Error(fmt.Sprintf("Metric type=%v is unsupported", metricType))
		ginContext.String(http.StatusBadRequest, "Metric type is unsupported")
	}

}

func (h *MetricHandler) handleCreateCounterMetric(ginContext *gin.Context) {
	metricName := ginContext.Param(constants.URLParamMetricName)
	metricValue := ginContext.Param(constants.URLParamMetricValue)

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

func (h *MetricHandler) handleCreateGaugeMetric(ginContext *gin.Context) {
	metricName := ginContext.Param(constants.URLParamMetricName)
	metricValue := ginContext.Param(constants.URLParamMetricValue)

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
