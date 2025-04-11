package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ruslanDantsov/osmetrics-server/internal/logging"
	"net/http"
)

type CommonHandler struct {
	Log logging.Logger
}

func NewCommonHandler(log logging.Logger) *CommonHandler {
	return &CommonHandler{
		Log: log,
	}
}

func (h *CommonHandler) ServeHTTP(ginContext *gin.Context) {
	h.Log.Error(fmt.Sprintf("Request is unsupported: url: %v; method: %v",
		ginContext.Request.RequestURI,
		ginContext.Request.Method))
	ginContext.String(http.StatusNotFound, "Request is unsupported")
}
