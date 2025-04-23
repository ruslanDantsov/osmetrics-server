package app

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/ruslanDantsov/osmetrics-server/internal/config"
	"github.com/ruslanDantsov/osmetrics-server/internal/service"
	"go.uber.org/zap"
	"log"
	"net/http"
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
	app.waitForServer()

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

func (app *AgentApp) waitForServer() {
	url := fmt.Sprintf("http://%v/health", app.config.Address)
	for {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			return
		}
		log.Println("Server not ready, waiting...")
		time.Sleep(2 * time.Second)
	}
}
