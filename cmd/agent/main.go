// Package main is the entry point for the osmetrics agent.
package main

import (
	"context"
	"fmt"
	"github.com/ruslanDantsov/osmetrics-server/internal/agent/app"
	"github.com/ruslanDantsov/osmetrics-server/internal/agent/config"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/logger"
	"io"
	"os"
	"os/signal"
	"syscall"
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

	defer logger.Log.Sync()

	logger.Log.Info("Starting agent...")

	agentApp := app.NewAgentApp(agentConfig, logger.Log)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := agentApp.Run(ctx); err != nil {
		logger.Log.Fatal(fmt.Sprintf("Agent start failed: %v", err.Error()))
	}
}
