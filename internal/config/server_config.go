package config

import (
	"flag"
	"fmt"
	"os"
)

type ServerConfig struct {
	Address string
}

func NewServerConfig() *ServerConfig {
	flagAddress := flag.String("a", "localhost:8080", "Address of the HTTP server")
	flag.Parse()

	if len(flag.Args()) > 0 {
		fmt.Fprintf(os.Stderr, "Error: unknown flags detected: %v\n", flag.Args())
		os.Exit(1)
	}
	address := getEnvOrDefault("ADDRESS", *flagAddress)

	return &ServerConfig{
		Address: address,
	}
}
