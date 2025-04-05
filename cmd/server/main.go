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
	gaugePostHandler := handler.NewGaugePostHandler(storage, log)
	gaugeGetHandler := handler.NewGaugeGetHandler(storage, log)
	counterPostHandler := handler.NewCounterPostHandler(storage, log)
	counterGetHandler := handler.NewCounterGetHandler(storage, log)
	commonHandler := handler.NewCommonHandler(log)

	startServer(commonHandler, gaugePostHandler, gaugeGetHandler, counterPostHandler, counterGetHandler)
}

func startServer(commonHandler *handler.CommonHandler,
	gaugePostHandler *handler.GaugePostHandler,
	gaugeGetHandler *handler.GaugeGetHandler,
	counterPostHandler *handler.CounterPostHandler,
	counterGetHandler *handler.CounterGetHandler) {

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, commonHandler.ServeHTTP)
	mux.HandleFunc("POST /update/gauge/{type}/{value}", gaugePostHandler.ServeHTTP)
	mux.HandleFunc("GET /gauge/{type}", gaugeGetHandler.ServeHTTP)
	mux.HandleFunc("POST /update/counter/{type}/{value}", counterPostHandler.ServeHTTP)
	mux.HandleFunc("GET /counter/{type}", counterGetHandler.ServeHTTP)

	err := http.ListenAndServe(`localhost:8080`, mux)
	if err != nil {
		panic(err)
	}
}
