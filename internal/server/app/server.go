package app

import (
	"context"
	"fmt"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/model"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/model/enum"
	"github.com/ruslanDantsov/osmetrics-server/internal/server/config"
	"github.com/ruslanDantsov/osmetrics-server/internal/server/handler"
	"github.com/ruslanDantsov/osmetrics-server/internal/server/handler/metric"
	"github.com/ruslanDantsov/osmetrics-server/internal/server/repository/file"
	"github.com/ruslanDantsov/osmetrics-server/internal/server/repository/memory"
	"github.com/ruslanDantsov/osmetrics-server/internal/server/repository/postgre"
	"github.com/ruslanDantsov/osmetrics-server/proto/metrics"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
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
	grpcHandler        *handler.ServerMetricsHandler
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

	serverMetricHandler := handler.NewServerMetricsHandler(storage, log)

	return &ServerApp{
		cfg:         cfg,
		logger:      log,
		grpcHandler: serverMetricHandler,
		storage:     storage,
	}, nil
}

// Run запускает HTTP-сервер и регистрирует все маршруты,
// настраивает middleware и маршруты профилирования pprof.

func (app *ServerApp) Run() error {
	listener, err := net.Listen("tcp", app.cfg.Address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", app.cfg.Address, err)
	}

	grpcServer := grpc.NewServer()
	metrics.RegisterMetricsServiceServer(grpcServer, app.grpcHandler)

	app.logger.Info("gRPC server started", zap.String("address", app.cfg.Address))
	return grpcServer.Serve(listener)
}

// Close завершает работу приложения.
func (app *ServerApp) Close() {
	app.storage.Close()
}
