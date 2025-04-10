package config

import (
	"github.com/alecthomas/kingpin/v2"
	"time"
)

type AgentConfig struct {
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

func NewAgentConfig(cliArgs []string) *AgentConfig {
	config := AgentConfig{}
	app := kingpin.New("agentApp", "Agent application")

	var pollIntervalSeconds int
	var reportIntervalSeconds int

	app.
		Flag("a", "Address of the HTTP server").
		Short('a').
		Envar("ADDRESS").
		Default("localhost:8080").
		StringVar(&config.Address)

	app.
		Flag("r", "Frequency (in seconds) for sending reports to the server").
		Short('r').
		Envar("REPORT_INTERVAL").
		Default("10").
		IntVar(&reportIntervalSeconds)

	app.
		Flag("p", "Frequency (in seconds) for polling metrics from runtime").
		Short('p').
		Envar("POLL_INTERVAL").
		Default("2").
		IntVar(&pollIntervalSeconds)

	_, err := app.Parse(cliArgs)
	if err != nil {
		panic(err)
	}

	config.ReportInterval = time.Duration(reportIntervalSeconds) * time.Second
	config.PollInterval = time.Duration(pollIntervalSeconds) * time.Second

	return &config
}
