package handler

import (
	"fmt"
	"github.com/ruslanDantsov/osmetrics-server/internal/logging"
	"github.com/ruslanDantsov/osmetrics-server/internal/model"
	"github.com/ruslanDantsov/osmetrics-server/internal/repository"
	"net/http"
)

type CounterPostHandler struct {
	Storage repository.Storager
	Log     logging.Logger
}

func NewCounterPostHandler(storage repository.Storager, log logging.Logger) *CounterPostHandler {
	return &CounterPostHandler{
		Storage: storage,
		Log:     log,
	}
}

func (h *CounterPostHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	h.Log.Info(fmt.Sprintf("Handle request %v", request.RequestURI))

	contentType := request.Header.Get("Content-Type")
	if contentType != "text/plain" {
		h.Log.Error("Bad Request: Content-Type must be text/plain")
		http.Error(response, "Bad Request: Content-Type must be text/plain", http.StatusBadRequest)
		return
	}

	counterType := request.PathValue("type")
	counterValue := request.PathValue("value")

	counterModel, err := model.NewCounterMetricModelWithRawTypes(counterType, counterValue)
	if err != nil {
		h.Log.Error(err.Error())
		http.Error(response, err.Error(), http.StatusBadRequest)
	}

	_, err = h.Storage.SaveCounterMetric(counterModel)
	if err != nil {
		h.Log.Error(err.Error())
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	response.WriteHeader(http.StatusOK)
}
