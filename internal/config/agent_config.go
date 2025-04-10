package config

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
	"time"
)

type AgentConfig struct {
	Address                 string        `long:"address" short:"a" env:"ADDRESS" default:"localhost:8080" description:"Address of the HTTP server"`
	ReportIntervalInSeconds int           `long:"report" short:"r" env:"REPORT_INTERVAL" default:"10" description:"Frequency (in seconds) for sending reports to the server"`
	PollIntervalInSeconds   int           `long:"poll" short:"p" env:"POLL_INTERVAL" default:"2" description:"Frequency (in seconds) for polling metrics from runtime"`
	ReportInterval          time.Duration `long:"-" description:"Derived duration from ReportIntervalInSeconds"`
	PollInterval            time.Duration `no:"-" description:"Derived duration from PollSeconds"`
}

func NewAgentConfig(cliArgs []string) *AgentConfig {
	config := &AgentConfig{}
	parser := flags.NewParser(config, flags.Default)

	_, err := parser.ParseArgs(cliArgs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Convert seconds to durations
	config.ReportInterval = time.Duration(config.ReportIntervalInSeconds) * time.Second
	config.PollInterval = time.Duration(config.PollIntervalInSeconds) * time.Second

	return config
}
