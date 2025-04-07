package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type AgentConfig struct {
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

func NewAgentConfig() *AgentConfig {
	flagAddress := flag.String("a", "localhost:8080", "Address of the HTTP server")
	flagReportInterval := flag.Int("r", 10, "Frequency (in seconds) for sending reports to the server")
	flagPollInterval := flag.Int("p", 2, "Frequency (in seconds) for polling metrics from runtime")

	flag.Parse()

	if len(flag.Args()) > 0 {
		fmt.Fprintf(os.Stderr, "Error: unknown flags detected: %v\n", flag.Args())
		os.Exit(1)
	}

	address := getEnvOrDefault("ADDRESS", *flagAddress)
	reportInterval := getEnvTimeOrDefault("REPORT_INTERVAL", time.Duration(*flagReportInterval)*time.Second)
	pollInterval := getEnvTimeOrDefault("POLL_INTERVAL", time.Duration(*flagPollInterval)*time.Second)

	return &AgentConfig{
		Address:        address,
		ReportInterval: reportInterval,
		PollInterval:   pollInterval,
	}
}

func getEnvOrDefault(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func getEnvTimeOrDefault(key string, defaultVal time.Duration) time.Duration {
	if valStr, ok := os.LookupEnv(key); ok {
		if val, err := strconv.Atoi(valStr); err == nil {
			return time.Duration(val) * time.Second
		}
	}
	return defaultVal
}
