package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/ruslanDantsov/osmetrics-server/internal/config"
	"github.com/ruslanDantsov/osmetrics-server/internal/logging"
	"github.com/ruslanDantsov/osmetrics-server/internal/service"
	"time"
)

func main() {
	agentConfig := config.NewAgentConfig()
	log := logging.NewStdoutLogger()
	client := resty.New()
	metricService := service.NewMetricService(log, client, agentConfig)
	go func() {
		ticker := time.NewTicker(agentConfig.ReportInterval)
		defer ticker.Stop()

		for range ticker.C {
			metricService.CollectMetrics()
		}
	}()

	go func() {
		ticker := time.NewTicker(agentConfig.ReportInterval)
		defer ticker.Stop()

		for range ticker.C {
			metricService.SendMetrics()
		}
	}()

	select {}
}
