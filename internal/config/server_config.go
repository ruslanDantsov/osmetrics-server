package config

import (
	"github.com/alecthomas/kingpin/v2"
)

type ServerConfig struct {
	Address string
}

func NewServerConfig(cliArgs []string) *ServerConfig {
	config := ServerConfig{}
	app := kingpin.New("serverApp", "Server application")
	app.
		Flag("host", "Server host address").
		Envar("ADDRESS").
		Default("localhost:8080").
		StringVar(&config.Address)

	_, err := app.Parse(cliArgs)
	if err != nil {
		panic(err)
	}

	return &config
}
