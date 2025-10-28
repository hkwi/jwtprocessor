package jwtprocessor

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"slices"

	jwt "github.com/golang-jwt/jwt/v5"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processorhelper"
)

type Config struct {
	Alg               string `mapstructure:"alg"`
	PrivateKey        string `mapstructure:"private_key"`        // Base64-encoded bytes for HMAC, PEM for RSA/ECDSA/Ed25519
	ResourceAttribute string `mapstructure:"resource_attribute"` // Attribute key to store the signed JWT
	ScopeAttribute    string `mapstructure:"scope_attribute"`
	Attribute         string `mapstructure:"attribute"`
}

func (c Config) getSigningKey() (any, error) {
	if !slices.Contains(jwt.GetAlgorithms(), c.Alg) {
		return nil, fmt.Errorf("unsupported algorithm: %s", c.Alg)
	}
	switch jwt.GetSigningMethod(c.Alg).(type) {
	case *jwt.SigningMethodRSA, *jwt.SigningMethodRSAPSS:
		return jwt.ParseRSAPrivateKeyFromPEM([]byte(c.PrivateKey))
	case *jwt.SigningMethodECDSA:
		return jwt.ParseECPrivateKeyFromPEM([]byte(c.PrivateKey))
	case *jwt.SigningMethodEd25519:
		return jwt.ParseEdPrivateKeyFromPEM([]byte(c.PrivateKey))
	case *jwt.SigningMethodHMAC:
		return base64.StdEncoding.DecodeString(c.PrivateKey)
	default:
		return nil, fmt.Errorf("unknown signing method: %s", c.Alg)
	}
}

func (c Config) Validate() error {
	if _, err := c.getSigningKey(); err != nil {
		return err
	}
	return nil
}

func sign(method jwt.SigningMethod, pkey any, body map[string]any, path, attr_name string) error {
	if len(attr_name) == 0 {
		return nil
	}
	tok := jwt.NewWithClaims(method, jwt.MapClaims(body))
	if signedToken, err := tok.SignedString(pkey); err != nil {
		return err
	} else {
		var info map[string]any
		if len(path) > 0 {
			info = body[path].(map[string]any)
		} else {
			info = body
		}
		var overwrite bool = false
		rattrs := []any{}
		if attrs, ok := info["attributes"]; ok {
			for _, attr := range attrs.([]any) {
				attr := attr.(map[string]any)
				if attr["key"] == attr_name {
					attr["value"] = map[string]string{"stringValue": signedToken}
					overwrite = true
				} else {
					rattrs = append(rattrs, attr)
				}
			}
		}
		if !overwrite {
			rattrs = append(rattrs, map[string]any{
				"key": attr_name,
				"value": map[string]string{
					"stringValue": signedToken,
				},
			})
		}
		info["attributes"] = rattrs
	}
	return nil
}

func NewFactory() processor.Factory {
	return processor.NewFactory(
		component.MustNewType("jwt"),
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
	pkey, err := pcfg.getSigningKey()
	if err != nil {
		return nil, err
	}
	return processorhelper.NewLogs(
		ctx,
		set,
		cfg,
		next,
		func(ctx context.Context, ld plog.Logs) (plog.Logs, error) {
			v := make(map[string]any, 1)
			m := &plog.JSONMarshaler{}
			if b, err := m.MarshalLogs(ld); err != nil {
				return ld, err
			} else if err := json.Unmarshal(b, &v); err != nil {
				return ld, err
			} else if resourceLogs, ok := v["resourceLogs"]; ok {
				for _, resourceLog := range resourceLogs.([]any) {
					resourceLog := resourceLog.(map[string]any)
					sign(alg, pkey, resourceLog, "resource", pcfg.ResourceAttribute)
					if scopeLogs, ok := resourceLog["scopeLogs"]; ok {
						for _, scopeLog := range scopeLogs.([]any) {
							scopeLog := scopeLog.(map[string]any)
							sign(alg, pkey, scopeLog, "scope", pcfg.ScopeAttribute)
							if logRecords, ok := scopeLog["logRecords"]; ok {
								for _, logRecord := range logRecords.([]any) {
									logRecord := logRecord.(map[string]any)
									sign(alg, pkey, logRecord, "", pcfg.Attribute)
								}
							}
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
			return ld, nil
		},
		processorhelper.WithCapabilities(consumer.Capabilities{MutatesData: true}),
	)
}

func createMetrics(ctx context.Context, set processor.Settings, cfg component.Config, next consumer.Metrics) (processor.Metrics, error) {
	pcfg := cfg.(*Config) // Ensure the config is of type *Config
	alg := jwt.GetSigningMethod(pcfg.Alg)
	pkey, err := pcfg.getSigningKey()
	if err != nil {
		return nil, err
	}
	return processorhelper.NewMetrics(
		ctx,
		set,
		cfg,
		next,
		func(ctx context.Context, md pmetric.Metrics) (pmetric.Metrics, error) {
			v := make(map[string]any, 1)
			m := &pmetric.JSONMarshaler{}
			if b, err := m.MarshalMetrics(md); err != nil {
				return md, err
			} else if err := json.Unmarshal(b, &v); err != nil {
				return md, err
			} else if resourceMetrics, ok := v["resourceMetrics"]; ok {
				for _, resourceMetric := range resourceMetrics.([]any) {
					resourceMetric := resourceMetric.(map[string]any)
					sign(alg, pkey, resourceMetric, "resource", pcfg.ResourceAttribute)
					if scopeMetrics, ok := resourceMetric["scopeMetrics"]; ok {
						for _, scopeMetric := range scopeMetrics.([]any) {
							scopeMetric := scopeMetric.(map[string]any)
							sign(alg, pkey, scopeMetric, "scope", pcfg.ScopeAttribute)
							if metrics, ok := scopeMetric["metrics"]; ok {
								for _, metric := range metrics.([]any) {
									metric := metric.(map[string]any)
									for _, dataPointType := range []string{"gauge", "sum", "histogram", "exponential_histogram", "summary"} {
										if dataPoints, ok := metric[dataPointType]; ok {
											for _, dataPoint := range dataPoints.([]any) {
												dataPoint := dataPoint.(map[string]any)
												sign(alg, pkey, dataPoint, "", pcfg.Attribute)
											}
										}
									}
								}
							}
						}
					}
				}
				if b, err := json.Marshal(v); err != nil {
					return md, err
				} else {
					un := &pmetric.JSONUnmarshaler{}
					return un.UnmarshalMetrics(b)
				}
			}
			return md, nil
		},
		processorhelper.WithCapabilities(consumer.Capabilities{MutatesData: true}),
	)
}

func createTrace(ctx context.Context, set processor.Settings, cfg component.Config, next consumer.Traces) (processor.Traces, error) {
	pcfg := cfg.(*Config) // Ensure the config is of type *Config
	alg := jwt.GetSigningMethod(pcfg.Alg)
	pkey, err := pcfg.getSigningKey()
	if err != nil {
		return nil, err
	}
	return processorhelper.NewTraces(
		ctx,
		set,
		cfg,
		next,
		func(ctx context.Context, td ptrace.Traces) (ptrace.Traces, error) {
			v := make(map[string]any, 1)
			m := &ptrace.JSONMarshaler{}
			if b, err := m.MarshalTraces(td); err != nil {
				return td, err
			} else if err := json.Unmarshal(b, &v); err != nil {
				return td, err
			} else if resourceSpans, ok := v["resourceSpans"]; ok {
				for _, resourceSpan := range resourceSpans.([]any) {
					resourceSpan := resourceSpan.(map[string]any)
					sign(alg, pkey, resourceSpan, "resource", pcfg.ResourceAttribute)
					if scopeSpans, ok := resourceSpan["scopeSpans"]; ok {
						for _, scopeSpan := range scopeSpans.([]any) {
							scopeSpan := scopeSpan.(map[string]any)
							sign(alg, pkey, scopeSpan, "scope", pcfg.ScopeAttribute)
							if spans, ok := scopeSpan["spans"]; ok {
								for _, span := range spans.([]any) {
									span := span.(map[string]any)
									sign(alg, pkey, span, "", pcfg.Attribute)
								}
							}
						}
					}
				}
				if b, err := json.Marshal(v); err != nil {
					return td, err
				} else {
					un := &ptrace.JSONUnmarshaler{}
					return un.UnmarshalTraces(b)
				}
			}
			return td, nil
		},
		processorhelper.WithCapabilities(consumer.Capabilities{MutatesData: true}),
	)
}
