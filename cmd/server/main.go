package main

import (
	"github.com/gin-gonic/gin"
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

	router := gin.Default()
	router.GET(`/`, metricListHandler.ServeHTTP)
	router.GET("/value/:type/:name", metricGetHandler.ServeHTTP)
	router.POST("/update/:type/:name/:value", metricPostHandler.ServeHTTP)
	router.Any(`/:path/`, commonHandler.ServeHTTP)

	err := http.ListenAndServe(`localhost:8080`, router)
	if err != nil {
		panic(err)
	}
}
