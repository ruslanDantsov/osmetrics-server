package config

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
	"time"
)

// AgentConfig содержит параметры конфигурации агента
type AgentConfig struct {
	// Address — адрес HTTP-сервера, к которому агент будет отправлять метрики.
	Address string `long:"address" short:"a" env:"ADDRESS" default:"localhost:8090" description:"Address of the HTTP server"`

	// ReportIntervalInSeconds — частота (в секундах), с которой агент отправляет отчёты на сервер.
	ReportIntervalInSeconds int `long:"report" short:"r" env:"REPORT_INTERVAL" default:"10" description:"Frequency (in seconds) for sending reports to the server"`

	// PollIntervalInSeconds — частота (в секундах), с которой агент опрашивает метрики.
	PollIntervalInSeconds int `long:"poll" short:"p" env:"POLL_INTERVAL" default:"2" description:"Frequency (in seconds) for polling metrics from runtime"`

	// ReportInterval — производное значение из ReportIntervalInSeconds в формате time.Duration.
	ReportInterval time.Duration `long:"-" description:"Derived duration from ReportIntervalInSeconds"`

	// PollInterval — производное значение из PollIntervalInSeconds в формате time.Duration.
	PollInterval time.Duration `no:"-" description:"Derived duration from PollSeconds"`

	// LogLevel — уровень логирования (например, DEBUG, INFO, WARN, ERROR).
	LogLevel string `short:"v" long:"log" env:"LOG_LEVEL" default:"INFO" description:"Log Level"`

	// HashKey — секретный ключ для хеширования метрик.
	HashKey string `long:"key" short:"k" env:"KEY" description:"Secret key for hashing"`

	// RateLimit — количество рабочих потоков, отправляющих метрики на сервер.
	RateLimit int `long:"rate" short:"l" env:"RATE_LIMIT" default:"2" description:"Count of workers for sending metrics to the server"`

	CryptoPubKeyPath string `short:"c" long:"crypto-key" env:"CRYPTO_KEY" description:"path to public key"`
}

// NewAgentConfig создаёт и инициализирует конфигурацию агента,
// используя переданные аргументы командной строки и/или переменных окружения.
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
