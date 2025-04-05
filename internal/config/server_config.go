package config

type Server struct {
	Port   string
	Domain string
}

type envVar struct {
	Key          string
	DefaultValue string
}

var (
	PortVar = envVar{
		Key:          "PORT",
		DefaultValue: "8080",
	}
	DomainVar = envVar{
		Key:          "DOMAIN",
		DefaultValue: "localhost",
	}
)

func NewServerConfig() *Server {
	return &Server{
		Port:   getEnv(PortVar),
		Domain: getEnv(DomainVar),
	}
}

func getEnv(variable envVar) string {
	if len(variable.Key) == 0 {
		panic("Unknown key")
	}
	return variable.DefaultValue
}
