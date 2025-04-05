package handler

import (
	"encoding/json"
	"fmt"
	"github.com/ruslanDantsov/osmetrics-server/internal/logging"
	"github.com/ruslanDantsov/osmetrics-server/internal/model/enum/metric"
	"github.com/ruslanDantsov/osmetrics-server/internal/repository"
	"net/http"
)

type CounterGetHandler struct {
	Storage repository.Storager
	Log     logging.Logger
}

func NewCounterGetHandler(storage repository.Storager, log logging.Logger) *CounterGetHandler {
	return &CounterGetHandler{
		Storage: storage,
		Log:     log,
	}
}

func (h *CounterGetHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	h.Log.Info(fmt.Sprintf("Handle GET request %v", request.RequestURI))

	rawMetricType := request.PathValue("type")

	metricType, err := metric.ParseMetricType(rawMetricType)
	if err != nil {
		h.Log.Error(err.Error())
		http.Error(response, "Metric not found", http.StatusNotFound)
		return
	}

	counterModel, found := h.Storage.GetCounterMetric(metricType)

	if !found {
		h.Log.Error(fmt.Sprintf("The counter_metric type=%v not found", metricType))
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
