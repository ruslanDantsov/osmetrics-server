package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ruslanDantsov/osmetrics-server/internal/config"
	"github.com/ruslanDantsov/osmetrics-server/internal/handler"
	"github.com/ruslanDantsov/osmetrics-server/internal/handler/metric"
	"github.com/ruslanDantsov/osmetrics-server/internal/middleware"
	"github.com/ruslanDantsov/osmetrics-server/internal/repository"
	"go.uber.org/zap"
	"net/http"
)

type ServerApp struct {
	config             *config.ServerConfig
	logger             *zap.Logger
	getMetricHandler   *metric.GetMetricHandler
	storeMetricHandler *metric.StoreMetricHandler
	commonHandler      *handler.CommonHandler
	healthHandler      *handler.HealthHandler
	dbHealthHandler    *handler.DbHandler
	db                 *pgxpool.Pool
}

func NewServerApp(cfg *config.ServerConfig, log *zap.Logger) (*ServerApp, error) {
	db, err := pgxpool.New(context.Background(), cfg.DatabaseConnection)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	err = db.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	baseStorage := repository.NewMemStorage(*log)
	persistentStorage := repository.NewPersistentStorage(baseStorage, cfg.FileStoragePath, cfg.StoreInterval, *log, cfg.Restore)
	getMetricHandler := metric.NewGetMetricHandler(persistentStorage, *log)
	storeMetricHandler := metric.NewStoreMetricHandler(persistentStorage, *log)
	commonHandler := handler.NewCommonHandler(*log)
	healthHandler := handler.NewHealthHandler(*log)
	dbHealthHandler := handler.NewDbHandler(*log, db)

	return &ServerApp{
		config:             cfg,
		logger:             log,
		getMetricHandler:   getMetricHandler,
		storeMetricHandler: storeMetricHandler,
		commonHandler:      commonHandler,
		healthHandler:      healthHandler,
		dbHealthHandler:    dbHealthHandler,
		db:                 db,
	}, nil
}

func (app *ServerApp) Run() error {
	router := gin.Default()

	router.Use(middleware.NewLoggerRequestMiddleware(app.logger))
	router.Use(middleware.NewGzipCompressionMiddleware())
	router.Use(middleware.NewGzipDecompressionMiddleware())

	router.GET(`/`, app.getMetricHandler.List)
	router.GET("/health", app.healthHandler.GetHealth)
	router.GET("/ping", app.dbHealthHandler.GetDbHealth)
	router.GET("/value/:type/:name", app.getMetricHandler.Get)
	router.POST("/value", app.getMetricHandler.GetJSON)
	router.POST("/update", app.storeMetricHandler.StoreJSON)
	router.POST("/update/:type/:name/:value", app.storeMetricHandler.Store)
	router.Any(`/:path/`, app.commonHandler.ServeHTTP)

	return http.ListenAndServe(app.config.Address, router)
}

func (app *ServerApp) Close() {
	if app.db != nil {
		app.db.Close()
	}
}
