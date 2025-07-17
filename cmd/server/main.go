package main

import (
	"github.com/ruslanDantsov/osmetrics-server/internal/pkg/shared/logger"
	"github.com/ruslanDantsov/osmetrics-server/internal/server/app"
	"github.com/ruslanDantsov/osmetrics-server/internal/server/config"
	"go.uber.org/zap"
	"os"
)

func main() {
	serverConfig, err := config.NewServerConfig(os.Args[1:])

	if err != nil {
		logger.Log.Fatal("Config initialized failed: %v", zap.Error(err))
	}

	if err := logger.Initialized(serverConfig.LogLevel); err != nil {
		logger.Log.Fatal("Logger initialized failed: %v", zap.Error(err))
	}

	defer logger.Log.Sync()
	defer logger.Log.Sync()

	logger.Log.Info("Starting server...")

	serverApp, err := app.NewServerApp(serverConfig, logger.Log)
	if err != nil {
		logger.Log.Fatal("Unable to config Server", zap.Error(err))
	}

	if err := serverApp.Run(); err != nil {
		logger.Log.Fatal("Server start failed: %v", zap.Error(err))
	}

	defer serverApp.Close()

}
