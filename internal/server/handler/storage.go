package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type StorageHealthChecker interface {
	HealthCheck(ctx context.Context) error
}

type DBHandler struct {
	Log zap.Logger
	db  StorageHealthChecker
}

func NewDBHandler(log zap.Logger, db StorageHealthChecker) *DBHandler {
	return &DBHandler{
		Log: log,
		db:  db,
	}
}

func (h *DBHandler) GetDBHealth(ginContext *gin.Context) {
	ctx := ginContext.Request.Context()
	if err := h.db.HealthCheck(ctx); err != nil {
		h.Log.Warn("DB health check failed", zap.Error(err))
		ginContext.String(http.StatusInternalServerError, "DB health check failed")
		return
	}

	h.Log.Info("DB health check passed")
	ginContext.String(http.StatusOK, "DB health check passed")
}
