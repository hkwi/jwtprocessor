
# jwtprocessor: OpenTelemetry Collector Processor for JWT Signing

`jwtprocessor` is a custom processor for the OpenTelemetry Collector that generates JWT signatures at multiple levels (resource, scope, and record/span/datapoint) and stores them in specified attributes. This enables tamper detection and trust for observability data.

---

## Features

- **Multi-level JWT signing**: Generates JWT signatures at the resource, scope, and record/span/datapoint levels, storing them in configurable attributes
- **Supported data types**: Logs, Metrics, Traces
- **Signing algorithms**: HMAC, RSA, ECDSA, Ed25519 (via [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt))
- **Collector Builder compatible**: Easily integrated with the OpenTelemetry Collector builder

---

## Usage

1. Add `jwtprocessor` to your `builder-config.yaml`:
	 ```yaml
	 processors:
		 - gomod: github.com/hkwi/jwtprocessor main
	 ```
2. Build with the [OpenTelemetry Collector builder](https://opentelemetry.io/docs/collector/custom-collector/):
	 ```bash
	 ocb --config builder-config.yaml
	 ```

---

## Configuration Example

```yaml
processors:
	jwt:
		alg: "EdDSA"
		private_key: |
			-----BEGIN PRIVATE KEY-----
			[base64 encoded key]
			-----END PRIVATE KEY-----
		resource_attribute: "signed.resource"
		scope_attribute: "signed.scope"
		attribute: "signed.record"
```

- `alg`: Signing algorithm (e.g., "EdDSA", "RS256", etc.)
- `private_key`: PEM or Base64 format (depending on algorithm)
- `resource_attribute`, `scope_attribute`, `attribute`: Attribute names to store the JWT

---

## How it works

- For each log, metric, or trace, the processor signs the content at the resource, scope, and record/span/datapoint levels (if configured), and stores the JWT in the specified attribute.
- For validation, extract the JWT from the attribute and compare its decoded content with the actual data.

---

## References

- [OpenTelemetry Collector custom processor development](https://opentelemetry.io/docs/collector/custom-collector/)
- [golang-jwt/jwt/v5 documentation](https://pkg.go.dev/github.com/golang-jwt/jwt/v5)

---

## License

Apache License 2.0

