package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// HealthHandler обрабатывает запросы, связанные с проверкой состояния (health check) сервиса.
type HealthHandler struct {
	Log zap.Logger
}

// NewHealthHandler создаёт новый экземпляр HealthHandler.
func NewHealthHandler(log zap.Logger) *HealthHandler {
	return &HealthHandler{
		Log: log,
	}
}

// GetHealth обрабатывает HTTP-запрос на проверку состояния сервиса.
func (h *HealthHandler) GetHealth(ginContext *gin.Context) {
	h.Log.Info("I am healthy")
	ginContext.String(http.StatusOK, "I am healthy")
}
