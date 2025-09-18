package app

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/crypto"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/model"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/model/enum"
	"github.com/ruslanDantsov/osmetrics-server/internal/server/config"
	"github.com/ruslanDantsov/osmetrics-server/internal/server/handler"
	"github.com/ruslanDantsov/osmetrics-server/internal/server/handler/metric"
	"github.com/ruslanDantsov/osmetrics-server/internal/server/middleware"
	"github.com/ruslanDantsov/osmetrics-server/internal/server/repository/file"
	"github.com/ruslanDantsov/osmetrics-server/internal/server/repository/memory"
	"github.com/ruslanDantsov/osmetrics-server/internal/server/repository/postgre"
	"net/http/pprof"

	"go.uber.org/zap"
	"net/http"
)

// ServerApp представляет основное приложение сервера.
// Оно хранит конфигурацию, логгер, хендлеры и хранилище данных.
type ServerApp struct {
	cfg                *config.ServerConfig
	logger             *zap.Logger
	getMetricHandler   *metric.GetMetricHandler
	storeMetricHandler *metric.StoreMetricHandler
	commonHandler      *handler.CommonHandler
	healthHandler      *handler.HealthHandler
	dbHealthHandler    *handler.DBHandler
	storage            Storager
}

// Storager определяет интерфейс для взаимодействия с хранилищем метрик.
type Storager interface {
	GetKnownMetrics(ctx context.Context) []string
	GetMetric(ctx context.Context, metricID enum.MetricID) (*model.Metrics, bool)
	SaveMetric(ctx context.Context, metric *model.Metrics) (*model.Metrics, error)
	SaveAllMetrics(ctx context.Context, metricList model.MetricsList) (model.MetricsList, error)
	HealthCheck(ctx context.Context) error
	Close()
}

// NewServerApp создаёт и инициализирует экземпляр ServerApp.
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

	getMetricHandler := metric.NewGetMetricHandler(storage, *log)
	storeMetricHandler := metric.NewStoreMetricHandler(storage, *log)
	commonHandler := handler.NewCommonHandler(*log)
	healthHandler := handler.NewHealthHandler(*log)

	dbHealthHandler := handler.NewDBHandler(*log, storage)

	return &ServerApp{
		cfg:                cfg,
		logger:             log,
		getMetricHandler:   getMetricHandler,
		storeMetricHandler: storeMetricHandler,
		commonHandler:      commonHandler,
		healthHandler:      healthHandler,
		dbHealthHandler:    dbHealthHandler,
		storage:            storage,
	}, nil
}

// Run запускает HTTP-сервер и регистрирует все маршруты,
// настраивает middleware и маршруты профилирования pprof.
func (app *ServerApp) Run() error {
	router := gin.Default()

	var decryptMW gin.HandlerFunc
	if len(app.cfg.CryptoPrivateKeyPath) > 0 {
		privKey, err := crypto.LoadPrivateKey(app.cfg.CryptoPrivateKeyPath)
		if err != nil {
			app.logger.Fatal("failed to load private key", zap.Error(err))
		}
		decryptMW = middleware.NewDecryptPayloadMiddleware(privKey, app.logger)
	}

	router.Use(middleware.NewLoggerRequestMiddleware(app.logger))
	if len(app.cfg.HashKey) != 0 {
		router.Use(middleware.HashCheckerMiddleware(app.cfg.HashKey, app.logger))
	}
	router.Use(middleware.NewGzipCompressionMiddleware())
	router.Use(middleware.NewGzipDecompressionMiddleware())
	router.Use(middleware.TrustedSubnetMiddleware(app.cfg.TrustedSubnet, app.logger))

	router.GET(`/`, app.getMetricHandler.List)
	router.GET("/health", app.healthHandler.GetHealth)
	router.GET("/ping", app.dbHealthHandler.GetDBHealth)
	router.GET("/value/:type/:name", app.getMetricHandler.Get)
	if decryptMW != nil {
		router.POST("/value/", decryptMW, app.getMetricHandler.GetJSON)
		router.POST("/update", decryptMW, app.storeMetricHandler.StoreJSON)
		router.POST("/updates/", decryptMW, app.storeMetricHandler.StoreBatchJSON)
		router.POST("/update/:type/:name/:value", decryptMW, app.storeMetricHandler.Store)
	} else {
		router.POST("/value/", app.getMetricHandler.GetJSON)
		router.POST("/update", app.storeMetricHandler.StoreJSON)
		router.POST("/updates/", app.storeMetricHandler.StoreBatchJSON)
		router.POST("/update/:type/:name/:value", app.storeMetricHandler.Store)
	}
	router.Any(`/:path`, app.commonHandler.ServeHTTP)

	pprofGroup := router.Group("/debug/pprof")
	{
		pprofGroup.GET("/", gin.WrapH(http.HandlerFunc(pprof.Index)))
		pprofGroup.GET("/cmdline", gin.WrapH(http.HandlerFunc(pprof.Cmdline)))
		pprofGroup.GET("/profile", gin.WrapH(http.HandlerFunc(pprof.Profile)))
		pprofGroup.GET("/symbol", gin.WrapH(http.HandlerFunc(pprof.Symbol)))
		pprofGroup.GET("/trace", gin.WrapH(http.HandlerFunc(pprof.Trace)))
		pprofGroup.GET("/heap", gin.WrapH(http.HandlerFunc(pprof.Handler("heap").ServeHTTP)))
		pprofGroup.GET("/goroutine", gin.WrapH(http.HandlerFunc(pprof.Handler("goroutine").ServeHTTP)))
		pprofGroup.GET("/threadcreate", gin.WrapH(http.HandlerFunc(pprof.Handler("threadcreate").ServeHTTP)))
		pprofGroup.GET("/block", gin.WrapH(http.HandlerFunc(pprof.Handler("block").ServeHTTP)))
		pprofGroup.GET("/mutex", gin.WrapH(http.HandlerFunc(pprof.Handler("mutex").ServeHTTP)))
	}

	return http.ListenAndServe(app.cfg.Address, router)
}

// Close завершает работу приложения.
func (app *ServerApp) Close() {
	app.storage.Close()
}
