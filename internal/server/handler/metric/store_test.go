package metric

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/model"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"strings"
)

// mockSaver реализует интерфейс MetricSaver.
type MockSaver struct{}

func (m *MockSaver) SaveMetric(_ context.Context, metric *model.Metrics) (*model.Metrics, error) {
	return metric, nil
}

func (m *MockSaver) SaveAllMetrics(_ context.Context, _ model.MetricsList) (model.MetricsList, error) {
	return nil, nil
}

func ExampleStoreMetricHandler_Store_gauge() {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := NewStoreMetricHandler(&MockSaver{}, *zaptest.NewLogger(nil))
	r.POST("/update/:type/:name/:value", handler.Store)

	req := httptest.NewRequest(http.MethodPost, "/update/gauge/Alloc/123.45", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	// Output:
	// 200
}

func ExampleStoreMetricHandler_Store_counter() {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := NewStoreMetricHandler(&MockSaver{}, *zaptest.NewLogger(nil))
	r.POST("/update/:type/:name/:value", handler.Store)

	req := httptest.NewRequest(http.MethodPost, "/update/counter/PollCount/42", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	// Output:
	// 200
}

func ExampleStoreMetricHandler_Store_invalidType() {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := NewStoreMetricHandler(&MockSaver{}, *zaptest.NewLogger(nil))
	r.POST("/update/:type/:name/:value", handler.Store)

	req := httptest.NewRequest(http.MethodPost, "/update/unknown/Any/42", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)
	fmt.Println(strings.TrimSpace(w.Body.String()))

	// Output:
	// 400
	// metric type unknown is unsupported
}
