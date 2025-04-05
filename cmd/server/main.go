package main

import (
	"github.com/ruslanDantsov/osmetrics-server/internal/handler"
	"github.com/ruslanDantsov/osmetrics-server/internal/logging"
	"github.com/ruslanDantsov/osmetrics-server/internal/repository"
	"net/http"
)

func main() {
	//TODO: move to config
	log := logging.NewStdoutLogger()
	storage := repository.NewMemStorage(log)
	metricGetHandler := handler.NewMetricGetHandler(storage, log)
	metricPostHandler := handler.NewMetricPostHandler(storage, log)
	commonHandler := handler.NewCommonHandler(log)

	startServer(commonHandler, metricPostHandler, metricGetHandler)
}

func startServer(commonHandler *handler.CommonHandler,
	metricPostHandler *handler.MetricPostHandler,
	metricGetHandler *handler.MetricGetHandler) {

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, commonHandler.ServeHTTP)
	mux.HandleFunc("GET /{type}/{name}", metricGetHandler.ServeHTTP)
	mux.HandleFunc("POST /update/{type}/{name}/{value}", metricPostHandler.ServeHTTP)

	err := http.ListenAndServe(`localhost:8080`, mux)
	if err != nil {
		panic(err)
	}
}
