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
	metricListHandler := handler.NewMetricListHandler(storage, log)

	startServer(commonHandler, metricPostHandler, metricGetHandler, metricListHandler)
}

func startServer(commonHandler *handler.CommonHandler,
	metricPostHandler *handler.MetricPostHandler,
	metricGetHandler *handler.MetricGetHandler,
	metricListHandler *handler.MetricListHandler) {

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, metricListHandler.ServeHTTP)
	mux.HandleFunc("GET /value/{type}/{name}", metricGetHandler.ServeHTTP)
	mux.HandleFunc("POST /update/{type}/{name}/{value}", metricPostHandler.ServeHTTP)
	mux.HandleFunc(`/{path}/`, commonHandler.ServeHTTP)

	err := http.ListenAndServe(`localhost:8080`, mux)
	if err != nil {
		panic(err)
	}
}
