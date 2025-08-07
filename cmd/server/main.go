// Package main is the entry point for the osmetrics server.
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

	serverConfig, err := config.NewServerConfig(os.Args[1:])

	if err != nil {
		logger.Log.Fatal("Config initialized failed: %v", zap.Error(err))
	}

	if err := logger.Initialized(serverConfig.LogLevel); err != nil {
		logger.Log.Fatal("Logger initialized failed: %v", zap.Error(err))
	}

	defer func(Log *zap.Logger) {
		err := Log.Sync()
		if err != nil {
			logger.Log.Error(err.Error())
		}
	}(logger.Log)

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
