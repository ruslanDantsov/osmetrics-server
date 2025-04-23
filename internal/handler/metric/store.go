package metric

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ruslanDantsov/osmetrics-server/internal/constants"
	"github.com/ruslanDantsov/osmetrics-server/internal/model"
	"net/http"
	"strings"
)

func (h *MetricHandler) Store(ginContext *gin.Context) {
	metricType := strings.ToLower(ginContext.Param(constants.URLParamMetricType))
	metricName := ginContext.Param(constants.URLParamMetricName)
	metricValue := ginContext.Param(constants.URLParamMetricValue)

	var metricRequest *model.Metrics
	var err error

	switch metricType {
	case "counter":
		metricRequest, err = model.NewMetricWithRawValues("Counter", metricName, metricValue)
	case "gauge":
		metricRequest, err = model.NewMetricWithRawValues("Gauge", metricName, metricValue)
	default:
		h.Log.Warn(fmt.Sprintf("Metric type=%v is unsupported", metricType))
		metricRequest, err = model.NewMetricWithRawValues(metricType, metricName, metricValue)
	}

	if err != nil {
		h.Log.Error(err.Error())
		ginContext.String(http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.Storage.SaveMetric(metricRequest)
	if err != nil {
		h.Log.Error(err.Error())
		ginContext.String(http.StatusBadRequest, err.Error())
		return
	}

	ginContext.Status(http.StatusOK)
}
