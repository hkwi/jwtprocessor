package jwtprocessor

import (
	"go.opentelemetry.io/collector/component"
)

type Config struct {
	component.ProcessorConfig
	PrivateKeyType string `mapstructure:"private_key_type"`
	PrivateKey     string `mapstructure:"private_key"`
}
