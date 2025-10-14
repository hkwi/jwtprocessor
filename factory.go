package jwtprocessor

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/processor"
)

const TypeStr = "jwtprocessor"

func NewFactory() component.ProcessorFactory {
	return processor.NewFactory(
		TypeStr,
		createDefaultConfig,
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		ProcessorConfig: component.NewProcessorConfig(TypeStr),
	}
}
