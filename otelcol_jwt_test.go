package jwtprocessor

import (
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestOtelcolJwtOutput(t *testing.T) {
	if err := os.Chdir("test"); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	// ../dist/otelcol-jwt --config otelcol-config.yaml をバックグラウンドで起動
	cmd := exec.Command("../dist/otelcol", "--config", "otelcol-config.yaml")
	logFile := "sample_out.log"
	f, err := os.Create(logFile)
	if err != nil {
		t.Fatalf("failed to create log file: %v", err)
	}
	defer f.Close()
	cmd.Stdout = f
	cmd.Stderr = f
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start otelcol-jwt: %v", err)
	}
	// プロセスはバックグラウンドで動作

	// sample_out.log の生成を待つ
	for i := 0; i < 2; i++ {
		if _, err := os.Stat(logFile); err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	// sample_out.log に signed. が含まれるか検査
	data, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}
	if !strings.Contains(string(data), "signed.") {
		// 2秒待って再検査
		time.Sleep(2 * time.Second)
		data, err = os.ReadFile(logFile)
		if err != nil {
			t.Fatalf("failed to read log file after wait: %v", err)
		}
		if !strings.Contains(string(data), "signed.") {
			t.Fatalf("'signed.' not found in log file")
		}
	}

	// プロセス終了
	cmd.Process.Kill()
}
