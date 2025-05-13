package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ruslanDantsov/osmetrics-server/internal/constants"
	"go.uber.org/zap"
	"net/http"
)

type DBHandler struct {
	Log zap.Logger
	db  *pgxpool.Pool
}

func NewDBHandler(log zap.Logger, db *pgxpool.Pool) *DBHandler {
	return &DBHandler{
		Log: log,
		db:  db,
	}
}

func (h *DBHandler) GetDBHealth(ginContext *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DBPingTimeout)
	defer cancel()

	if err := h.db.Ping(ctx); err != nil {
		h.Log.Warn("DB health check failed", zap.Error(err))
		ginContext.String(http.StatusInternalServerError, "DB health check failed")
		return
	}

	h.Log.Info("DB health check passed")
	ginContext.String(http.StatusOK, "DB health check passed")
}
