package config

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

type ServerConfig struct {
	Address string `short:"a" long:"address" env:"ADDRESS" default:"localhost:8080" description:"Server host address"`
}

func NewServerConfig(cliArgs []string) *ServerConfig {
	fmt.Println("Start getting data for server config")

	config := &ServerConfig{}
	parser := flags.NewParser(config, flags.Default)

	_, err := parser.ParseArgs(cliArgs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("The data for the server has been loaded")
	return config
}
