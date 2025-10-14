# Research: opentelemetry-collector-processor

## Task 1: Processor Scaffolding Best Practices
- **Decision**: Scaffold new processor by copying the structure of the dnslookup processor, following OpenTelemetry Collector conventions.
- **Rationale**: Ensures compatibility and maintainability; leverages proven patterns.
- **Alternatives considered**: Creating processor from scratch (risk of missing conventions), using other processor as reference (dnslookup is closest in function).

## Task 2: builder-config.yaml for Custom Processor
- **Decision**: Add jwt processor to builder-config.yaml under processors, and enable otlp receiver, otlp exporter, debug exporter, kafka exporter, and filter processor.
- **Rationale**: Follows collector builder documentation and ensures all required components are enabled for testing.
- **Alternatives considered**: Manual config editing (risk of errors), omitting components (would not meet requirements).

## Task 3: Collector Build and Runtime Validation
- **Decision**: Validate by building collector with builder-config.yaml and running it to confirm jwt processor is recognized and active.
- **Rationale**: Ensures integration is successful and processor is ready for further development.
- **Alternatives considered**: Only static build validation (does not confirm runtime behavior).

---

All clarifications resolved. Ready for Phase 1 design.
