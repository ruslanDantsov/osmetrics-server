package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestConfig_FromEnv(t *testing.T) {
	expectedAddress := "localhost:1234"

	os.Setenv("ADDRESS", expectedAddress)
	defer os.Unsetenv("ADDRESS") // clean up

	config := NewServerConfig([]string{})

	assert.Equal(t, expectedAddress, config.Address)

}

func TestConfig_FromDefaultVlue(t *testing.T) {
	expectedAddress := "localhost:8080"

	config := NewServerConfig([]string{})

	assert.Equal(t, expectedAddress, config.Address)

}
