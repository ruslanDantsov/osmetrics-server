package metric

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/model"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/model/enum"
	"go.uber.org/zap/zaptest"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

// mockStorage реализует интерфейс MetricGetter.
type MockStorage struct{}

func (m *MockStorage) GetMetric(_ context.Context, metricID enum.MetricID) (*model.Metrics, bool) {
	switch metricID {
	case "Alloc":
		v := 123.45
		return &model.Metrics{
			ID:    "Alloc",
			MType: "gauge",
			Value: &v,
		}, true
	case "PollCount":
		d := int64(42)
		return &model.Metrics{
			ID:    "PollCount",
			MType: "counter",
			Delta: &d,
		}, true
	default:
		return nil, false
	}
}

func (m *MockStorage) GetKnownMetrics(_ context.Context) []string {
	return []string{"Alloc", "PollCount"}
}

func ExampleGetMetricHandler_Get_gauge() {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := NewGetMetricHandler(&MockStorage{}, *zaptest.NewLogger(nil))
	r.GET("/value/:type/:name", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/value/gauge/Alloc", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	body, _ := io.ReadAll(w.Body)
	fmt.Println(w.Code)
	fmt.Println(strings.TrimSpace(string(body)))
}

func ExampleGetMetricHandler_Get_counter() {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := NewGetMetricHandler(&MockStorage{}, *zaptest.NewLogger(nil))
	r.GET("/value/:type/:name", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/value/counter/PollCount", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	body, _ := io.ReadAll(w.Body)
	fmt.Println(w.Code)
	fmt.Println(strings.TrimSpace(string(body)))
}
