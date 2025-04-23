package metric

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"
	"github.com/ruslanDantsov/osmetrics-server/internal/model"
	"github.com/ruslanDantsov/osmetrics-server/internal/model/enum"
	"net/http"
)

func (h *MetricHandler) GetJson(ginContext *gin.Context) {
	var metricRequest model.Metrics

	if err := easyjson.UnmarshalFromReader(ginContext.Request.Body, &metricRequest); err != nil {
		h.Log.Error(fmt.Sprintf("Error on unmarshal data from request. %v", err))
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	switch metricRequest.MType {
	case "gauge":
		h.handleGetGaugeMetric(ginContext)
	case "counter":
		h.handleGetCounterMetric1(ginContext)
	default:
		h.Log.Error(fmt.Sprintf("Metric type=%v is unsupported", metricRequest.MType))
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Metric type is unsupported", "type": metricRequest.MType})
	}
}

func (h *MetricHandler) handleGetCounterMetric1(ginContext *gin.Context) {
	//rawMetricName := ginContext.Param(constants.URLParamMetricName)

	delta := int64(14)
	resp := model.Metrics{
		ID:    enum.MetricId("Alloc1"),
		MType: "Counter",
		Delta: &delta,
	}

	//counterModel, found := h.Storage.GetCounterMetric(metricName)
	//
	//if !found {
	//	h.Log.Error(fmt.Sprintf("The counter_metric name=%v not found", metricName))
	//	ginContext.String(http.StatusNotFound, "Metric not found")
	//	return
	//}

	ginContext.Header("Content-Type", "application/json")
	ginContext.Writer.WriteHeader(http.StatusOK)

	_, err := easyjson.MarshalToWriter(resp, ginContext.Writer)
	if err != nil {
		h.Log.Error(fmt.Sprintf("Error on marshal metric data. %v", err))
		ginContext.String(http.StatusNotFound, "Metric not found")
	}
}

//func (h *MetricHandler) handleGetGaugeMetric(ginContext *gin.Context) {
//	rawMetricName := ginContext.Param(constants.URLParamMetricName)
//
//	metricName, err := metric.ParseMetricName(rawMetricName)
//	if err != nil {
//		h.Log.Error(err.Error())
//		ginContext.String(http.StatusNotFound, "Metric name is unsupported")
//		return
//	}
//
//	gaugeModel, found := h.Storage.GetGaugeMetric(metricName)
//
//	if !found {
//		h.Log.Error(fmt.Sprintf("The gauge_metric name=%v not found", metricName))
//		ginContext.String(http.StatusNotFound, "Metric not found")
//		return
//	}
//
//	ginContext.Header("Content-Type", "text/html")
//	ginContext.String(http.StatusOK, strconv.FormatFloat(gaugeModel.Value, 'f', -1, 64))
//}
