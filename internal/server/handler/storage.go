package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// StorageHealthChecker описывает интерфейс для проверки состояния хранилища (базы данных).
type StorageHealthChecker interface {
	HealthCheck(ctx context.Context) error
}

// DBHandler обрабатывает HTTP-запросы, связанные с проверкой состояния базы данных.
type DBHandler struct {
	Log zap.Logger
	db  StorageHealthChecker
}

// NewDBHandler создает новый экземпляр DBHandler.
func NewDBHandler(log zap.Logger, db StorageHealthChecker) *DBHandler {
	return &DBHandler{
		Log: log,
		db:  db,
	}
}

// GetDBHealth обрабатывает HTTP-запрос проверки здоровья базы данных.
// В случае ошибки возвращает HTTP 500, иначе HTTP 200 с соответствующим сообщением.
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
