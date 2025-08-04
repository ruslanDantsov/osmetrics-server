package main

import (
	"context"
	"fmt"
	"github.com/ruslanDantsov/osmetrics-server/internal/agent/app"
	"github.com/ruslanDantsov/osmetrics-server/internal/agent/config"
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
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
