package handler

import (
	"fmt"
	"github.com/ruslanDantsov/osmetrics-server/internal/logging"
	"github.com/ruslanDantsov/osmetrics-server/internal/model"
	"github.com/ruslanDantsov/osmetrics-server/internal/repository"
	"net/http"
)

type GaugePostHandler struct {
	Storage repository.Storager
	Log     logging.Logger
}

func NewGaugePostHandler(storage repository.Storager, log logging.Logger) *GaugePostHandler {
	return &GaugePostHandler{
		Storage: storage,
		Log:     log,
	}
}

func (h *GaugePostHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	h.Log.Info(fmt.Sprintf("Handle POST request %v", request.RequestURI))

	contentType := request.Header.Get("Content-Type")
	if contentType != "text/plain" {
		h.Log.Error("Bad Request: Content-Type must be text/plain")
		http.Error(response, "Bad Request: Content-Type must be text/plain", http.StatusBadRequest)
		return
	}

	gaugeType := request.PathValue("type")
	gaugeValue := request.PathValue("value")

	gaugeModel, err := model.NewGaugeMetricModelWithRawTypes(gaugeType, gaugeValue)
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
