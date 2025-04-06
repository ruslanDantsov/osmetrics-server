package config

import (
	"flag"
	"time"
)

type AgentConfig struct {
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

func NewAgentConfig() *AgentConfig {
	addr := flag.String("a", "localhost:8080", "Address of the HTTP server")
	reportInterval := flag.Int("r", 10, "Frequency (in seconds) for sending reports to the server")
	pollInterval := flag.Int("p", 2, "Frequency (in seconds) for polling metrics from runtime")

	flag.Parse()

	return &AgentConfig{
		Address:        *addr,
		ReportInterval: time.Duration(*reportInterval) * time.Second,
		PollInterval:   time.Duration(*pollInterval) * time.Second,
	}
}
