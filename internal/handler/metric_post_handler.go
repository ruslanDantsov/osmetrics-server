package handler

import (
	"fmt"
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

func (h *MetricPostHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	h.Log.Info(fmt.Sprintf("Handle request %v", request.RequestURI))

	contentType := request.Header.Get("Content-Type")
	if contentType != "text/plain" {
		h.Log.Error("Content-Type must be text/plain")
		http.Error(response, "Content-Type must be text/plain", http.StatusBadRequest)
		return
	}

	metricType := request.PathValue("type")
	switch metricType {
	case "gauge":
		h.handlePostGaugeMetric(response, request)
	case "counter":
		h.handlePostCounterMetric(response, request)
	default:
		h.Log.Error(fmt.Sprintf("Metric type=%v is unsupported", metricType))
		http.Error(response, "Metric type is unsupported", http.StatusBadRequest)
	}

}

func (h *MetricPostHandler) handlePostCounterMetric(response http.ResponseWriter, request *http.Request) {
	metricName := request.PathValue("name")
	metricValue := request.PathValue("value")

	counterModel, err := model.NewCounterMetricModelWithRawValues(metricName, metricValue)
	if err != nil {
		h.Log.Error(err.Error())
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.Storage.SaveCounterMetric(counterModel)
	if err != nil {
		h.Log.Error(err.Error())
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	response.WriteHeader(http.StatusOK)
}

func (h *MetricPostHandler) handlePostGaugeMetric(response http.ResponseWriter, request *http.Request) {

	metricName := request.PathValue("name")
	metricValue := request.PathValue("value")

	gaugeModel, err := model.NewGaugeMetricModelWithRawValues(metricName, metricValue)
	if err != nil {
		h.Log.Error(err.Error())
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.Storage.SaveGaugeMetric(gaugeModel)
	if err != nil {
		h.Log.Error(err.Error())
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	response.WriteHeader(http.StatusOK)
}
