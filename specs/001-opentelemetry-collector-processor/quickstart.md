# Quickstart: opentelemetry-collector-processor

## 1. Clone the repository and switch to feature branch
```bash
git clone [REPO_URL]
cd jwtprocessor
git checkout 001-opentelemetry-collector-processor
```

## 2. Scaffold the JWT processor
- Copy the dnslookup processor directory as a template
- Rename and update files to create jwtprocessor skeleton

## 3. Update builder-config.yaml
- Add the following components:
  - Receivers: otlp
  - Exporters: otlp, debug, kafka
  - Processors: filter, jwtprocessor
- Ensure all components are enabled in the config

## 4. Build the collector
```bash
[collector-builder-command] --config builder-config.yaml
```

## 5. Validate the build
- Run the collector and confirm jwtprocessor is recognized and active

## 6. Next steps
- Implement JWT processing logic in jwtprocessor
- Add tests and documentation
