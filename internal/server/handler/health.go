package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type HealthHandler struct {
	Log zap.Logger
}

func NewHealthHandler(log zap.Logger) *HealthHandler {
	return &HealthHandler{
		Log: log,
	}
}

func (h *HealthHandler) GetHealth(ginContext *gin.Context) {
	h.Log.Info("I am healthy")
	ginContext.String(http.StatusOK, "I am healthy")
}
