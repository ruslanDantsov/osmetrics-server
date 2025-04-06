package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/ruslanDantsov/osmetrics-server/internal/logging"
	"github.com/ruslanDantsov/osmetrics-server/internal/service"
	"time"
)

const (
	CollectorTimeInSecond time.Duration = 2
	SenderTimeInSecond    time.Duration = 10
)

func main() {
	log := logging.NewStdoutLogger()
	client := resty.New()
	metricService := service.NewMetricService(log, client)
	go func() {
		ticker := time.NewTicker(CollectorTimeInSecond * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			metricService.CollectMetrics()
		}
	}()

	go func() {
		ticker := time.NewTicker(SenderTimeInSecond * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			metricService.SendMetrics()
		}
	}()

	select {}
}
