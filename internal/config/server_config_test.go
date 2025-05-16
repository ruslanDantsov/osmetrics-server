package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestServerConfig_FromEnv(t *testing.T) {
	expectedAddress := "localhost:1234"

	os.Setenv("ADDRESS", expectedAddress)
	defer os.Unsetenv("ADDRESS")

	config, _ := NewServerConfig([]string{})

	assert.Equal(t, expectedAddress, config.Address)
}

func TestServerConfig_FromDefaultValue(t *testing.T) {
	expectedAddress := "localhost:8080"

	config, _ := NewServerConfig([]string{})

	assert.Equal(t, expectedAddress, config.Address)
}

func TestServerConfig_FromCommandLineArg(t *testing.T) {
	expectedAddress := "localhost:1234"

	config, _ := NewServerConfig([]string{"-a=localhost:1234"})

	assert.Equal(t, expectedAddress, config.Address)

}

func TestServerConfig_RestoreArgument_FromCommandLineArg(t *testing.T) {
	expectedRestore := true

	config, _ := NewServerConfig([]string{"-r=true"})

	assert.Equal(t, expectedRestore, config.Restore)

}
