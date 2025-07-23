package metric

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/model"
	"github.com/ruslanDantsov/osmetrics-server/internal/server/constants"
	"go.uber.org/zap/zaptest"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

func ExampleGetMetricHandler_GetJSON_gauge() {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := NewGetMetricHandler(&MockStorage{}, *zaptest.NewLogger(nil))
	r.POST("/value/", handler.GetJSON)

	metricReq := model.Metrics{
		ID:    "Alloc",
		MType: constants.GaugeMetricType,
	}
	var buf bytes.Buffer
	_, _ = easyjson.MarshalToWriter(&metricReq, &buf)

	req := httptest.NewRequest(http.MethodPost, "/value/", &buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	body, _ := io.ReadAll(w.Body)
	fmt.Println(w.Code)
	fmt.Println(strings.Contains(string(body), `"Alloc"`))

	// Output:
	// 200
	// true
}

func ExampleGetMetricHandler_GetJSON_counter() {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := NewGetMetricHandler(&MockStorage{}, *zaptest.NewLogger(nil))
	r.POST("/value/", handler.GetJSON)

	metricReq := model.Metrics{
		ID:    "PollCount",
		MType: constants.CounterMetricType,
	}
	var buf bytes.Buffer
	_, _ = easyjson.MarshalToWriter(&metricReq, &buf)

	req := httptest.NewRequest(http.MethodPost, "/value/", &buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	body, _ := io.ReadAll(w.Body)
	fmt.Println(w.Code)
	fmt.Println(strings.Contains(string(body), `"PollCount"`))

	// Output:
	// 200
	// true
}

func ExampleGetMetricHandler_GetJSON_notFound() {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := NewGetMetricHandler(&MockStorage{}, *zaptest.NewLogger(nil))
	r.POST("/value/", handler.GetJSON)

	metricReq := model.Metrics{
		ID:    "UnknownMetric",
		MType: constants.GaugeMetricType,
	}
	var buf bytes.Buffer
	_, _ = easyjson.MarshalToWriter(&metricReq, &buf)

	req := httptest.NewRequest(http.MethodPost, "/value/", &buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	body, _ := io.ReadAll(w.Body)
	fmt.Println(w.Code)
	fmt.Println(strings.Contains(string(body), `"Metric not found"`))

	// Output:
	// 404
	// true
}

func ExampleGetMetricHandler_GetJSON_badJSON() {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := NewGetMetricHandler(&MockStorage{}, *zaptest.NewLogger(nil))
	r.POST("/value/", handler.GetJSON)

	badBody := strings.NewReader(`{ this is invalid json }`)

	req := httptest.NewRequest(http.MethodPost, "/value/", badBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	body, _ := io.ReadAll(w.Body)
	fmt.Println(w.Code)
	fmt.Println(strings.Contains(string(body), `"invalid JSON"`))

	// Output:
	// 400
	// true
}
