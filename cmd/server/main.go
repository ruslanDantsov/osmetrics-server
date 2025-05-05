package main

import (
	"fmt"
	"github.com/ruslanDantsov/osmetrics-server/internal/app"
	"github.com/ruslanDantsov/osmetrics-server/internal/config"
	"github.com/ruslanDantsov/osmetrics-server/internal/logger"
	"os"
)

func main() {
	serverConfig, err := config.NewServerConfig(os.Args[1:])

	if err != nil {
		logger.Log.Fatal(fmt.Sprintf("Config initialized failed: %v", err.Error()))
	}

	if err := logger.Initialized(serverConfig.LogLevel); err != nil {
		logger.Log.Fatal(fmt.Sprintf("Logger initialized failed: %v", err.Error()))
	}

	defer logger.Log.Sync()

	logger.Log.Info("Starting server...")

	serverApp := app.NewServerApp(serverConfig, logger.Log)
	if err := serverApp.Run(); err != nil {
		logger.Log.Fatal(fmt.Sprintf("Server start failed: %v", err.Error()))
	}
}
