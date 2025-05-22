package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
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
	expectedReportInterval := 20 * time.Second

	config := NewAgentConfig([]string{"-r", "20"})

	assert.Equal(t, expectedReportInterval, config.ReportInterval)

}
