package handler

import (
	"encoding/json"
	"fmt"
	"github.com/ruslanDantsov/osmetrics-server/internal/logging"
	"github.com/ruslanDantsov/osmetrics-server/internal/model/enum/metric"
	"github.com/ruslanDantsov/osmetrics-server/internal/repository"
	"net/http"
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

func (h *MetricGetHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	h.Log.Info(fmt.Sprintf("Handle request %v", request.RequestURI))

	metricType := request.PathValue("type")
	switch metricType {
	case "gauge":
		h.handleGetGaugeMetric(response, request)
	case "counter":
		h.handleGetCounterMetric(response, request)
	default:
		h.Log.Error(fmt.Sprintf("Metric type=%v is unsupported", metricType))
		http.Error(response, "Metric type is unsupported", http.StatusBadRequest)
	}
}

func (h *MetricGetHandler) handleGetCounterMetric(response http.ResponseWriter, request *http.Request) {
	rawMetricName := request.PathValue("name")

	metricName, err := metric.ParseMetricName(rawMetricName)
	if err != nil {
		h.Log.Error(err.Error())
		http.Error(response, "Unsupported metric type", http.StatusNotFound)
		return
	}

	counterModel, found := h.Storage.GetCounterMetric(metricName)

	if !found {
		h.Log.Error(fmt.Sprintf("The counter_metric name=%v not found", metricName))
		http.Error(response, "Metric not found", http.StatusNotFound)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	if jsonErr := json.NewEncoder(response).Encode(counterModel); jsonErr != nil {
		h.Log.Error(jsonErr.Error())
		http.Error(response, jsonErr.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *MetricGetHandler) handleGetGaugeMetric(response http.ResponseWriter, request *http.Request) {
	rawMetricName := request.PathValue("name")

	metricName, err := metric.ParseMetricName(rawMetricName)
	if err != nil {
		h.Log.Error(err.Error())
		http.Error(response, "Unsupported metric type", http.StatusNotFound)
		return
	}

	gaugeModel, found := h.Storage.GetGaugeMetric(metricName)

	if !found {
		h.Log.Error(fmt.Sprintf("The gauge_metric name=%v not found", metricName))
		http.Error(response, "Metric not found", http.StatusNotFound)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	if jsonErr := json.NewEncoder(response).Encode(gaugeModel); jsonErr != nil {
		h.Log.Error(jsonErr.Error())
		http.Error(response, jsonErr.Error(), http.StatusInternalServerError)
		return
	}
}
