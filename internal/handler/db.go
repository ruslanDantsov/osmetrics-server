package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type StorageHealthChecker interface {
	Ping() error
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
	if err := h.db.Ping(); err != nil {
		h.Log.Warn("DB health check failed", zap.Error(err))
		ginContext.String(http.StatusInternalServerError, "DB health check failed")
		return
	}

	h.Log.Info("DB health check passed")
	ginContext.String(http.StatusOK, "DB health check passed")
}
