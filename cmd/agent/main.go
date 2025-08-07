// Package main is the entry point for the osmetrics agent.
package main

import (
	"fmt"
	"github.com/ruslanDantsov/osmetrics-server/internal/app"
	"github.com/ruslanDantsov/osmetrics-server/internal/config"
	"github.com/ruslanDantsov/osmetrics-server/internal/logger"
	"go.uber.org/zap"
	"io"
	"os"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func printBuildInfo(w io.Writer) {
	_, err := fmt.Fprintf(w, "Build version: %s\n", buildVersion)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "Build date: %s\n", buildDate)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "Build commit: %s\n", buildCommit)
	if err != nil {
		return
	}
}

func main() {
	printBuildInfo(os.Stdout)

	agentConfig := config.NewAgentConfig(os.Args[1:])

	if err := logger.Initialized(agentConfig.LogLevel); err != nil {
		logger.Log.Fatal(fmt.Sprintf("Logger initialized failed: %v", err.Error()))
	}

	defer func(Log *zap.Logger) {
		err := Log.Sync()
		if err != nil {
			logger.Log.Error(err.Error())
		}
	}(logger.Log)

	logger.Log.Info("Starting agent...")

	agentApp := app.NewAgentApp(agentConfig, logger.Log)
	if err := agentApp.Run(); err != nil {
		logger.Log.Fatal(fmt.Sprintf("Agent start failed: %v", err.Error()))
	}
}
