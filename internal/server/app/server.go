package app

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/model"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/model/enum"
	"github.com/ruslanDantsov/osmetrics-server/internal/server/config"
	handler2 "github.com/ruslanDantsov/osmetrics-server/internal/server/handler"
	metric2 "github.com/ruslanDantsov/osmetrics-server/internal/server/handler/metric"
	middleware2 "github.com/ruslanDantsov/osmetrics-server/internal/server/middleware"
	"github.com/ruslanDantsov/osmetrics-server/internal/server/repository/file"
	"github.com/ruslanDantsov/osmetrics-server/internal/server/repository/memory"
	"github.com/ruslanDantsov/osmetrics-server/internal/server/repository/postgre"

	"go.uber.org/zap"
	"net/http"
)

type ServerApp struct {
	config             *config.ServerConfig
	logger             *zap.Logger
	getMetricHandler   *metric2.GetMetricHandler
	storeMetricHandler *metric2.StoreMetricHandler
	commonHandler      *handler2.CommonHandler
	healthHandler      *handler2.HealthHandler
	dbHealthHandler    *handler2.DBHandler
	storage            Storager
}

type Storager interface {
	GetKnownMetrics(ctx context.Context) []string
	GetMetric(ctx context.Context, metricID enum.MetricID) (*model.Metrics, bool)
	SaveMetric(ctx context.Context, metric *model.Metrics) (*model.Metrics, error)
	SaveAllMetrics(ctx context.Context, metricList model.MetricsList) (model.MetricsList, error)
	HealthCheck(ctx context.Context) error
	Close()
}

func NewServerApp(cfg *config.ServerConfig, log *zap.Logger) (*ServerApp, error) {
	var storage Storager
	var err error

	if cfg.DatabaseConnection != "" {
		storage, err = postgre.NewPostgreStorage(*log, cfg.DatabaseConnection)
		if err != nil {
			return nil, err
		}
	} else {
		baseStorage := memory.NewMemStorage(*log)
		storage = file.NewPersistentStorage(baseStorage, cfg.FileStoragePath, cfg.StoreInterval, *log, cfg.Restore)
	}

	getMetricHandler := metric2.NewGetMetricHandler(storage, *log)
	storeMetricHandler := metric2.NewStoreMetricHandler(storage, *log)
	commonHandler := handler2.NewCommonHandler(*log)
	healthHandler := handler2.NewHealthHandler(*log)

	dbHealthHandler := handler2.NewDBHandler(*log, storage)

	return &ServerApp{
		config:             cfg,
		logger:             log,
		getMetricHandler:   getMetricHandler,
		storeMetricHandler: storeMetricHandler,
		commonHandler:      commonHandler,
		healthHandler:      healthHandler,
		dbHealthHandler:    dbHealthHandler,
		storage:            storage,
	}, nil
}

func (app *ServerApp) Run() error {
	router := gin.Default()

	router.Use(middleware2.NewLoggerRequestMiddleware(app.logger))
	router.Use(middleware2.NewGzipCompressionMiddleware())
	router.Use(middleware2.NewGzipDecompressionMiddleware())

	router.GET(`/`, app.getMetricHandler.List)
	router.GET("/health", app.healthHandler.GetHealth)
	router.GET("/ping", app.dbHealthHandler.GetDBHealth)
	router.GET("/value/:type/:name", app.getMetricHandler.Get)
	router.POST("/value", app.getMetricHandler.GetJSON)
	router.POST("/update", app.storeMetricHandler.StoreJSON)
	router.POST("/updates", app.storeMetricHandler.StoreBatchJSON)
	router.POST("/update/:type/:name/:value", app.storeMetricHandler.Store)
	router.Any(`/:path/`, app.commonHandler.ServeHTTP)

	return http.ListenAndServe(app.config.Address, router)
}

func (app *ServerApp) Close() {
	app.storage.Close()
}
