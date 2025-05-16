package config

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

type ServerConfig struct {
	Address                string        `short:"a" long:"address" env:"ADDRESS" default:"localhost:8080" description:"Server host address"`
	LogLevel               string        `short:"l" long:"log" env:"LOG_LEVEL" default:"INFO" description:"Log Level"`
	StoreIntervalInSeconds int           `short:"i" long:"interval" env:"STORE_INTERVAL" default:"300" description:"Interval in seconds for storing metrics to file"`
	StoreInterval          time.Duration `long:"-" description:"Derived duration from ReportIntervalInSeconds"`
	FileStoragePath        string        `short:"f" long:"path" env:"FILE_STORAGE_PATH" default:"" description:"Path to file with metrics data"`
	RestoreRaw             string        `short:"r" long:"restore" env:"RESTORE" description:"Flag indicating whether to load previously saved metrics data" no-ini:"true"`
	Restore                bool          `ignored:"true"`
	DatabaseConnection     string        `short:"d" long:"database" env:"DATABASE_DSN" description:"Database connection string"`
}

func NewServerConfig(cliArgs []string) (*ServerConfig, error) {
	config := &ServerConfig{}
	parser := flags.NewParser(config, flags.Default)

	_, err := parser.ParseArgs(cliArgs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if config.FileStoragePath == "" {
		config.FileStoragePath = getDefaultStoragePath()
	}

	if err := createStorageDirectory(config.FileStoragePath); err != nil {
		return nil, fmt.Errorf("failed to prepare storage directory: %w", err)
	}

	config.StoreInterval = time.Duration(config.StoreIntervalInSeconds) * time.Second

	if config.RestoreRaw != "" {
		val, err := strconv.ParseBool(config.RestoreRaw)
		if err != nil {
			return nil, fmt.Errorf("invalid value for --restore: %v", err)
		}
		config.Restore = val
	}

	return config, nil
}

func getDefaultStoragePath() string {
	const appName = "osmetrics-server"
	const fileName = "metrics.json"

	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData == "" {
			appData = os.Getenv("LOCALAPPDATA")
		}
		return filepath.Join(appData, appName, fileName)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, "."+appName, fileName)
}

func createStorageDirectory(filePath string) error {
	dir := filepath.Dir(filePath)

	if _, err := os.Stat(dir); err == nil {
		return nil
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("could not create storage directory %s: %w", dir, err)
	}

	return nil
}
