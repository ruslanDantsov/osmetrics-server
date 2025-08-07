package metric

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type KnownMetricsGetter interface {
	GetKnownMetrics(ctx context.Context) []string
}

func (h *GetMetricHandler) List(ginContext *gin.Context) {
	ctx := ginContext.Request.Context()

	var metricNames = h.Storage.GetKnownMetrics(ctx)
	htmlContent := "<html><head><title>Список метрик</title></head><body>"
	htmlContent += "<h1>List of known metrics:</h1>"
	htmlContent += "<ul>"

	for _, metricName := range metricNames {
		htmlContent += fmt.Sprintf("<li>%s</li>", metricName)
	}

	htmlContent += "</ul>"
	htmlContent += "</body></html>"

	ginContext.Header("Content-Type", "text/html")
	ginContext.Status(http.StatusOK)
	_, err := ginContext.Writer.Write([]byte(htmlContent))
	if err != nil {
		return
	}
}
