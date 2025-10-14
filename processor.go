package jwtprocessor

import (
	"context"

	"go.opentelemetry.io/collector/pdata/plog"
)

type jwtProcessor struct {
	config *Config
}

func newProcessor(cfg *Config) *jwtProcessor {
	return &jwtProcessor{config: cfg}
}

func (p *jwtProcessor) ProcessLogs(ctx context.Context, ld plog.Logs) (plog.Logs, error) {
	// TODO: JWT処理ロジック
	return ld, nil
}
