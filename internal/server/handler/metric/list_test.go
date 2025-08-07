package metric

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zaptest"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

func ExampleGetMetricHandler_List() {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := NewGetMetricHandler(&MockStorage{}, *zaptest.NewLogger(nil))
	r.GET("/", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	body, _ := io.ReadAll(w.Body)
	fmt.Println(w.Code)
	fmt.Println(strings.Contains(string(body), "<li>Alloc</li>"))
	fmt.Println(strings.Contains(string(body), "<li>PollCount</li>"))

	// Output:
	// 200
	// true
	// true
}
