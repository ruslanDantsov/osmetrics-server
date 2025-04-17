package config

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

type ServerConfig struct {
	Address  string `short:"a" long:"address" env:"ADDRESS" default:"localhost:8080" description:"Server host address"`
	LogLevel string `short:"l" long:"log" env:"LOG_LEVEL" default:"INFO" description:"Log Level"`
}

func NewServerConfig(cliArgs []string) *ServerConfig {
	config := &ServerConfig{}
	parser := flags.NewParser(config, flags.Default)

	_, err := parser.ParseArgs(cliArgs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return config
}
