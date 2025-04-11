package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
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

func (h *MetricListHandler) ServeHTTP(c *gin.Context) {
	h.Log.Info(fmt.Sprintf("Handle request %v", c.Request.RequestURI))
	var metricNames = h.Storage.GetKnownMetrics()
	htmlContent := "<html><head><title>Список метрик</title></head><body>"
	htmlContent += "<h1>List of known metrics:</h1>"
	htmlContent += "<ul>"

	for _, metricName := range metricNames {
		htmlContent += fmt.Sprintf("<li>%s</li>", metricName)
	}

	htmlContent += "</ul>"
	htmlContent += "</body></html>"

	c.Header("Content-Type", "text/html")
	c.Status(http.StatusOK)
	c.Writer.Write([]byte(htmlContent))
}
