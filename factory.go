package jwtprocessor

import (
	"context"
	"encoding/json"
	"fmt"

	jwt "github.com/golang-jwt/jwt/v5"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processorhelper"
)

type Config struct {
	Alg        string `mapstructure:"alg"`
	PrivateKey string `mapstructure:"private_key_pem"`
}

func getSigningMethod(alg, privateKey string) (any, error) {
	if key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey)); err == nil {
		return key, nil
	}
	if key, err := jwt.ParseEdPrivateKeyFromPEM([]byte(privateKey)); err == nil {
		return key, nil
	}
	if key, err := jwt.ParseECPrivateKeyFromPEM([]byte(privateKey)); err == nil {
		return key, nil
	}
	return nil, fmt.Errorf("Unknown PEM block format")
}

func (c *Config) Validate() error {
	// TODO: Alg と PrivateKey を合わせて検査したい。
	// HMAC support を追加してもいい
	if _, err := getSigningMethod(c.Alg, c.PrivateKey); err != nil {
		return err
	}
	return nil
}

func NewFactory() processor.Factory {
	return processor.NewFactory(
		component.MustNewType("jwtprocessor"),
		func() component.Config {
			return &Config{}
		},
		processor.WithLogs(createLog, component.StabilityLevelAlpha),
		processor.WithMetrics(createMetrics, component.StabilityLevelAlpha),
		processor.WithTraces(createTrace, component.StabilityLevelAlpha),
	)
}

func createLog(ctx context.Context, set processor.Settings, cfg component.Config, next consumer.Logs) (processor.Logs, error) {
	pcfg := cfg.(*Config) // Ensure the config is of type *Config
	alg := jwt.GetSigningMethod(pcfg.Alg)
	pkey, err := getSigningMethod(pcfg.Alg, pcfg.PrivateKey)
	if err != nil {
		return nil, err
	}
	return processorhelper.NewLogs(
		ctx,
		set,
		cfg,
		next,
		func(ctx context.Context, ld plog.Logs) (plog.Logs, error) {
			var v map[string]any
			m := &plog.JSONMarshaler{}
			if b, err := m.MarshalLogs(ld); err != nil {
				return ld, err
			} else if err := json.Unmarshal(b, &v); err != nil {
				return ld, err
			} else if resourceLogs, ok := v["resourceLogs"]; !ok {
				return ld, fmt.Errorf("no resourceLogs")
			} else {
				for _, resourceLog := range resourceLogs.([]any) {
					rsrc := resourceLog.(map[string]any)["resource"].(map[string]any)
					claim := jwt.MapClaims{
						"scopeLogs": resourceLog.(map[string]any)["scopeLogs"],
					}
					tok := jwt.NewWithClaims(alg, claim)
					if signedToken, err := tok.SignedString(pkey); err != nil {
						return ld, err
					} else {
						signed := map[string]any{
							"key": "signed.scopeLogs",
							"value": map[string]string{
								"stringValue": signedToken,
							},
						}
						if attributes, ok := rsrc["attributes"]; ok {
							rsrc["attributes"] = append(attributes.([]any), signed)
						} else {
							rsrc["attributes"] = []any{signed}
						}
					}
				}
				if b, err := json.Marshal(v); err != nil {
					return ld, err
				} else {
					un := &plog.JSONUnmarshaler{}
					return un.UnmarshalLogs(b)
				}
			}
		},
		processorhelper.WithCapabilities(consumer.Capabilities{MutatesData: true}),
	)
}

func createMetrics(ctx context.Context, set processor.Settings, cfg component.Config, next consumer.Metrics) (processor.Metrics, error) {
	return nil, nil
}

func createTrace(ctx context.Context, set processor.Settings, cfg component.Config, next consumer.Traces) (processor.Traces, error) {
	return nil, nil
}
