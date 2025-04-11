package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ruslanDantsov/osmetrics-server/internal/config"
	"github.com/ruslanDantsov/osmetrics-server/internal/handler"
	"github.com/ruslanDantsov/osmetrics-server/internal/logging"
	"github.com/ruslanDantsov/osmetrics-server/internal/repository"
	"net/http"
	"os"
)

func main() {
	log := logging.NewStdoutLogger()
	storage := repository.NewMemStorage(log)
	metricGetHandler := handler.NewMetricGetHandler(storage, log)
	metricPostHandler := handler.NewMetricPostHandler(storage, log)
	commonHandler := handler.NewCommonHandler(log)
	metricListHandler := handler.NewMetricListHandler(storage, log)
	serverConfig := config.NewServerConfig(os.Args[1:])

	startServer(commonHandler, metricPostHandler, metricGetHandler, metricListHandler, serverConfig)
}

func startServer(commonHandler *handler.CommonHandler,
	metricPostHandler *handler.MetricPostHandler,
	metricGetHandler *handler.MetricGetHandler,
	metricListHandler *handler.MetricListHandler,
	serverConfig *config.ServerConfig) {

	router := gin.Default()
	router.GET(`/`, metricListHandler.ServeHTTP)
	router.GET("/value/:type/:name", metricGetHandler.ServeHTTP)
	router.POST("/update/:type/:name/:value", metricPostHandler.ServeHTTP)
	router.Any(`/:path/`, commonHandler.ServeHTTP)

	err := http.ListenAndServe(serverConfig.Address, router)
	if err != nil {
		panic(err)
	}
}
