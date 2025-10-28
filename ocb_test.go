package jwtprocessor

import (
	"os"
	"os/exec"
	"testing"
)

func TestOCBWithConfig(t *testing.T) {
	cmd := exec.Command("ocb", "--config", "builder-config.yaml")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("ocb failed: %v\nOutput: %s", err, output)
	}
	if _, err := os.Stat("dist/otelcol"); err != nil {
		// the path is defined in builder-config.yaml
		t.Fatalf("dist/otelcol not found: %v", err)
	}
}
