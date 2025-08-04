package metric

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
)

func ExampleStoreMetricHandler_StoreJSON() {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := NewStoreMetricHandler(&MockSaver{}, *zaptest.NewLogger(nil))
	r.POST("/update/", handler.StoreJSON)

	// JSON тела запроса
	body := `{
        "id": "Alloc",
        "type": "gauge",
        "value": 123.45
    }`

	req := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	// Output:
	// 200
}

func ExampleStoreMetricHandler_StoreJSON_invalidJSON() {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := NewStoreMetricHandler(&MockSaver{}, *zaptest.NewLogger(nil))
	r.POST("/update/", handler.StoreJSON)

	req := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewBufferString("not valid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	// Output:
	// 400
}

func ExampleStoreMetricHandler_StoreBatchJSON() {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := NewStoreMetricHandler(&MockSaver{}, *zaptest.NewLogger(nil))
	r.POST("/updates/", handler.StoreBatchJSON)

	body := `[
        {
            "id": "Alloc",
            "type": "gauge",
            "value": 123.45
        },
        {
            "id": "PollCount",
            "type": "counter",
            "delta": 5
        }
    ]`

	req := httptest.NewRequest(http.MethodPost, "/updates/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	// Output:
	// 200
}

func ExampleStoreMetricHandler_StoreBatchJSON_invalidBody() {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := NewStoreMetricHandler(&MockSaver{}, *zaptest.NewLogger(nil))
	r.POST("/updates/", handler.StoreBatchJSON)

	req := httptest.NewRequest(http.MethodPost, "/updates/", bytes.NewBufferString("not an array"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	// Output:
	// 400
}
