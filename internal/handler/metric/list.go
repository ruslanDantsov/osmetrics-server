package metric

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type KnownMetricsGetter interface {
	GetKnownMetrics() []string
}

func (h *GetMetricHandler) List(c *gin.Context) {
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
