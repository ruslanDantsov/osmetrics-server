package app

import (
	"github.com/go-resty/resty/v2"
	"github.com/ruslanDantsov/osmetrics-server/internal/config"
	"github.com/ruslanDantsov/osmetrics-server/internal/service"
	"go.uber.org/zap"
	"sync"
	"time"
)

type AgentApp struct {
	config        *config.AgentConfig
	logger        *zap.Logger
	client        resty.Client
	metricService *service.MetricService
}

func NewAgentApp(cfg *config.AgentConfig, log *zap.Logger) *AgentApp {

	client := resty.New()
	metricService := service.NewMetricService(log, client, cfg)

	return &AgentApp{
		config:        cfg,
		logger:        log,
		client:        *client,
		metricService: metricService,
	}
}

func (app *AgentApp) Run() error {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		ticker := time.NewTicker(app.config.PollInterval)
		defer ticker.Stop()

		for range ticker.C {
			app.metricService.CollectMetrics()
		}
	}()

	go func() {
		wg.Done()
		ticker := time.NewTicker(app.config.ReportInterval)
		defer ticker.Stop()

		for range ticker.C {
			app.metricService.SendMetrics()
		}
	}()

	wg.Wait()

	return nil
}
