package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// CommonHandler предоставляет обработчик для неподдерживаемых HTTP-запросов.
type CommonHandler struct {
	Log zap.Logger
}

// NewCommonHandler создаёт и возвращает новый экземпляр CommonHandler.
func NewCommonHandler(log zap.Logger) *CommonHandler {
	return &CommonHandler{
		Log: log,
	}
}

// ServeHTTP обрабатывает неподдерживаемые HTTP-запросы.
//
// Он логирует информацию о запросе и возвращает клиенту статус 404.
func (h *CommonHandler) ServeHTTP(ginContext *gin.Context) {
	h.Log.Error(fmt.Sprintf("Request is unsupported: url: %v; method: %v",
		ginContext.Request.RequestURI,
		ginContext.Request.Method))
	ginContext.String(http.StatusNotFound, "Request is unsupported")
}
