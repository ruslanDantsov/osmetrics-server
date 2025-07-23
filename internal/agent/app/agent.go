package app

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/ruslanDantsov/osmetrics-server/internal/agent/config"
	"github.com/ruslanDantsov/osmetrics-server/internal/agent/constants"
	middleware2 "github.com/ruslanDantsov/osmetrics-server/internal/agent/middleware"
	"github.com/ruslanDantsov/osmetrics-server/internal/agent/service"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/model"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

// AgentApp представляет основное приложение агента, которое отвечает за сбор и передачу метрик.
// AgentApp хранить информацию для запуска агента
type AgentApp struct {
	config        *config.AgentConfig
	logger        *zap.Logger
	client        resty.Client
	metricService *service.MetricService
}

// NewAgentApp создает новый экземпляр AgentApp с заданной конфигурацией и логгером.
func NewAgentApp(cfg *config.AgentConfig, log *zap.Logger) *AgentApp {

	client := resty.New()
	client.OnBeforeRequest(middleware2.GzipRestyMiddleware())
	if len(cfg.HashKey) > 0 {
		client.OnBeforeRequest(middleware2.HashBodyRestyMiddleware(cfg.HashKey))
	}

	metricService := service.NewMetricService(log, client, cfg)

	return &AgentApp{
		config:        cfg,
		logger:        log,
		client:        *client,
		metricService: metricService,
	}
}

// Run запускает агент: собирает метрики, отправляет их и управляет жизненным циклом воркеров.
//
// Ожидает готовности сервера перед запуском сбора метрик.
// Использует контекст для управления завершением работы.
func (app *AgentApp) Run(ctx context.Context) error {
	if err := app.waitForServer(); err != nil {
		return err
	}

	metricChan := make(chan model.Metrics, constants.MetricChannelSize)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		ticker := time.NewTicker(app.config.PollInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				app.metricService.CollectMetrics(metricChan)
			case <-ctx.Done():
				app.logger.Info("Collector received shutdown signal")
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(app.config.PollInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				app.metricService.CollectAdditionalMetrics(metricChan)
			case <-ctx.Done():
				app.logger.Info("Additional collector received shutdown signal")
				return
			}
		}
	}()

	for i := 0; i < app.config.RateLimit; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			app.metricService.Worker(ctx, metricChan)
		}()
	}

	wg.Wait()
	close(metricChan)

	return nil
}

func (app *AgentApp) waitForServer() error {
	url := fmt.Sprintf(constants.ServerHealthCheckURL, app.config.Address)

	sleepDelay := 1 * time.Second
	attemps := 0
	for {
		resp, err := app.client.R().Get(url)

		if err == nil && resp.StatusCode() == http.StatusOK {
			app.logger.Info("Server is ready!")
			return nil
		}

		app.logger.Info("Server not ready, waiting...")
		time.Sleep(sleepDelay)

		sleepDelay += constants.IncreaseDelayForWaitingServer
		attemps++

		if sleepDelay > constants.MaxDelayForWaitingServer {
			break
		}

	}
	return fmt.Errorf("server didn't become ready after %v attempts", attemps)
}
