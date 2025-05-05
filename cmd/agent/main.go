package main

import (
	"fmt"
	"github.com/ruslanDantsov/osmetrics-server/internal/app"
	"github.com/ruslanDantsov/osmetrics-server/internal/config"
	"github.com/ruslanDantsov/osmetrics-server/internal/logger"
	"os"
)

func main() {
	agentConfig := config.NewAgentConfig(os.Args[1:])

	if err := logger.Initialized(agentConfig.LogLevel); err != nil {
		logger.Log.Fatal(fmt.Sprintf("Logger initialized failed: %v", err.Error()))
	}

	defer logger.Log.Sync()

	logger.Log.Info("Starting agent...")

	agentApp := app.NewAgentApp(agentConfig, logger.Log)
	if err := agentApp.Run(); err != nil {
		logger.Log.Fatal(fmt.Sprintf("Agent start failed: %v", err.Error()))
	}
}
