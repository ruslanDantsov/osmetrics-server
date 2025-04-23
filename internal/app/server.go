package app

import (
	"github.com/gin-gonic/gin"
	"github.com/ruslanDantsov/osmetrics-server/internal/config"
	"github.com/ruslanDantsov/osmetrics-server/internal/handler"
	"github.com/ruslanDantsov/osmetrics-server/internal/handler/metric"
	"github.com/ruslanDantsov/osmetrics-server/internal/middleware"
	"github.com/ruslanDantsov/osmetrics-server/internal/repository"
	"go.uber.org/zap"
	"net/http"
)

type ServerApp struct {
	config        *config.ServerConfig
	logger        *zap.Logger
	storage       repository.Storager
	metricHandler *metric.MetricHandler
	commonHandler *handler.CommonHandler
	healthHandler *handler.HealthHandler
}

func NewServerApp(cfg *config.ServerConfig, log *zap.Logger) *ServerApp {
	storage := repository.NewMemStorage(*log)
	metricHandler := metric.NewMetricHandler(storage, *log)
	commonHandler := handler.NewCommonHandler(*log)
	healthHandler := handler.NewHealthHandler(*log)

	return &ServerApp{
		config:        cfg,
		logger:        log,
		storage:       storage,
		metricHandler: metricHandler,
		commonHandler: commonHandler,
		healthHandler: healthHandler,
	}
}

func (app *ServerApp) Run() error {
	router := gin.Default()

	router.Use(middleware.NewLoggerRequestMiddleware(app.logger))

	router.GET(`/`, app.metricHandler.List)
	router.GET("/health", app.healthHandler.GetHealth)
	router.GET("/value/:type/:name", app.metricHandler.Get)
	router.POST("/value", app.metricHandler.GetJSON)
	router.POST("/update", app.metricHandler.StoreJSON)
	router.POST("/update/:type/:name/:value", app.metricHandler.Store)
	router.Any(`/:path/`, app.commonHandler.ServeHTTP)

	return http.ListenAndServe(app.config.Address, router)
}
