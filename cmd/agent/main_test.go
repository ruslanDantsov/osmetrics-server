package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestWriteBuildInfo(t *testing.T) {
	buildVersion = "1.0.0-agent-test"
	buildDate = "2025-08-07"
	buildCommit = "100agenttest"

	var buf bytes.Buffer
	printBuildInfo(&buf)

	output := buf.String()
	expected := fmt.Sprintf(
		"Build version: %s\nBuild date: %s\nBuild commit: %s\n",
		buildVersion, buildDate, buildCommit,
	)

	if !strings.Contains(output, expected) {
		t.Errorf("unexpected output:\ngot:\n%s\nwant:\n%s", output, expected)
	}
}
