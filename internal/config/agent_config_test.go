package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestAgentConfig_FromEnv(t *testing.T) {
	expectedAddress := "localhost:1234"

	os.Setenv("ADDRESS", expectedAddress)
	defer os.Unsetenv("ADDRESS")

	config := NewAgentConfig([]string{})

	assert.Equal(t, expectedAddress, config.Address)

}

func TestAgentConfig_FromDefaultValue(t *testing.T) {
	expectedAddress := "localhost:8080"

	config := NewAgentConfig([]string{})

	assert.Equal(t, expectedAddress, config.Address)

}

func TestAgentConfig_FromCommandLineArg(t *testing.T) {
	expectedAddress := "localhost:1234"

	config := NewAgentConfig([]string{"-a", "localhost:1234"})

	assert.Equal(t, expectedAddress, config.Address)

}
