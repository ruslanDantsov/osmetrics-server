package handler

import (
	"fmt"
	"github.com/ruslanDantsov/osmetrics-server/internal/logging"
	"github.com/ruslanDantsov/osmetrics-server/internal/repository"
	"net/http"
)

type MetricListHandler struct {
	Storage repository.Storager
	Log     logging.Logger
}

func NewMetricListHandler(storage repository.Storager, log logging.Logger) *MetricListHandler {
	return &MetricListHandler{
		Log:     log,
		Storage: storage,
	}
}

func (h *MetricListHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	h.Log.Info(fmt.Sprintf("Handle request %v", request.RequestURI))
	var metricNames = h.Storage.GetKnownMetrics()
	htmlContent := "<html><head><title>Список метрик</title></head><body>"
	htmlContent += "<h1>List of known metrics:</h1>"
	htmlContent += "<ul>"

	for _, metricName := range metricNames {
		htmlContent += fmt.Sprintf("<li>%s</li>", metricName)
	}

	htmlContent += "</ul>"
	htmlContent += "</body></html>"

	response.Header().Set("Content-Type", "text/html")
	response.WriteHeader(http.StatusOK)
	response.Write([]byte(htmlContent))
}
