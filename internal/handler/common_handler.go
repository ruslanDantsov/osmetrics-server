package handler

import (
	"fmt"
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

func (h *CommonHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	h.Log.Error(fmt.Sprintf("Request is unsupported: url: %v; method: %v", request.RequestURI, request.Method))
	http.Error(response, "Request is unsupported", http.StatusNotFound)
}
