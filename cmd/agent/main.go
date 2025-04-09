package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/ruslanDantsov/osmetrics-server/internal/config"
	"github.com/ruslanDantsov/osmetrics-server/internal/logging"
	"github.com/ruslanDantsov/osmetrics-server/internal/service"
	"sync"
	"time"
)

func main() {
	agentConfig := config.NewAgentConfig()
	log := logging.NewStdoutLogger()
	client := resty.New()
	metricService := service.NewMetricService(log, client, agentConfig)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		ticker := time.NewTicker(agentConfig.ReportInterval)
		defer ticker.Stop()

		for range ticker.C {
			metricService.CollectMetrics()
		}
	}()

	go func() {
		wg.Done()
		ticker := time.NewTicker(agentConfig.ReportInterval)
		defer ticker.Stop()

		for range ticker.C {
			metricService.SendMetrics()
		}
	}()

	wg.Wait()
}
