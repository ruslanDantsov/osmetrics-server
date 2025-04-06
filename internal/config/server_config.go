package config

import "flag"

type ServerConfig struct {
	Address string
}

func NewServerConfig() *ServerConfig {
	addr := flag.String("a", "localhost:8080", "Address of the HTTP server")
	flag.Parse()

	return &ServerConfig{
		Address: *addr,
	}
}
