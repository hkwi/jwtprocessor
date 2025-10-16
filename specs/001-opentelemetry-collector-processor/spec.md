# Feature Specification: opentelemetry-collector-processor

**Feature Branch**: `001-opentelemetry-collector-processor`  
**Created**: October 14, 2025  
**Status**: Draft  
**Input**: User description: "opentelemetry collector processorを作成します。まずはスケルトンであるopentelemetry-collector dnslookup processorを元にスケルトンを作成します。次にopentelemetry collector builderでビルドするためのbuilder-config.yamlを作成します。このbuilder-config.yamlでは、otelp receiver, otelp exporter, debug exporter, kafka exporter, filter processor とこの jwt processorを有効にします。opentelemetry collector builderでビルドできることを確認します。"

## User Scenarios & Testing *(mandatory)*

<!--
  IMPORTANT: User stories should be PRIORITIZED as user journeys ordered by importance.
  Each user story/journey must be INDEPENDENTLY TESTABLE - meaning if you implement just ONE of them,
  you should still have a viable MVP (Minimum Viable Product) that delivers value.
  
  Assign priorities (P1, P2, P3, etc.) to each story, where P1 is the most critical.
  Think of each story as a standalone slice of functionality that can be:
  - Developed independently
  - Tested independently
  - Deployed independently
  - Demonstrated to users independently
-->

### User Story 1 - Create JWT Processor Skeleton (Priority: P1)

As a developer, I want to create a new opentelemetry collector processor (jwt processor) skeleton based on the existing dnslookup processor, so that I can extend and customize it for JWT processing.

**Why this priority**: This is the foundational step for enabling JWT processing in the collector, and all subsequent work depends on having a working skeleton.

**Independent Test**: Can be fully tested by verifying that the new processor skeleton is created, compiles, and is recognized by the collector builder.

**Acceptance Scenarios**:

1. **Given** the repo and the dnslookup processor as reference, **When** the jwt processor skeleton is created, **Then** the code structure matches collector processor conventions and compiles without errors.
2. **Given** the new processor skeleton, **When** included in the collector builder config, **Then** the collector can be built successfully.

---

### User Story 2 - Collector Build and Artifact Verification (Priority: P2)

As a developer, I want to build the collector with builder-config.yaml (jwtprocessor有効化)し、Goテストで `ocb --config builder-config.yaml` を実行して `dist/otelcol-jwt` の生成を自動検証したい。

**Why this priority**: ビルド成果物の自動検証により、CI/CDや開発効率が向上する。

**Independent Test**: Goテストで `ocb` 実行と成果物の存在確認ができる。

**Acceptance Scenarios**:

1. **Given** builder-config.yamlでjwtprocessorを有効化し、**When** `ocb --config builder-config.yaml` をGoテストで実行、**Then** `dist/otelcol-jwt` が生成されていることを自動検証できる。

---

### User Story 3 - Collector実行とログ検証 (Priority: P3)

As a developer, I want to `dist/otelcol-jwt --config otelcol-config.yaml` を `test` ディレクトリでバックグラウンド起動し、`sample_out.log` の生成と `signed.` の有無をGoテストで自動検証したい。

**Why this priority**: 実行・ログ検証の自動化により、jwtprocessorの動作確認が容易になる。

**Independent Test**: Goテストで成果物の実行・ログ内容検証ができる。

**Acceptance Scenarios**:

1. **Given** `dist/otelcol-jwt` と `otelcol-config.yaml`、**When** `dist/otelcol-jwt --config otelcol-config.yaml` を `test` ディレクトリでバックグラウンド起動、**Then** `sample_out.log` が生成され、`signed.` を含むことを自動検証できる。

---

### Edge Cases

- What happens if the jwt processor skeleton has missing or invalid configuration?
- How does the collector builder handle a processor that fails to compile?
- What if the builder-config.yaml omits a required component?

### Edge Cases

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right edge cases.
-->

- What happens when [boundary condition]?
- How does system handle [error scenario]?

## Requirements *(mandatory)*

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right functional requirements.
-->

### Functional Requirements

- **FR-001**: System MUST provide a new opentelemetry collector processor skeleton named "jwt processor" based on the dnslookup processor structure.
- **FR-002**: System MUST ensure the jwt processor skeleton compiles and is recognized by the collector builder.
- **FR-003**: System MUST provide a builder-config.yaml that enables otlp receiver, otlp exporter, debug exporter, kafka exporter, filter processor, and the jwt processor.
- **FR-004**: System MUST build the collector successfully with the jwt processor and all specified components enabled.
- **FR-005**: System MUST validate that the built collector runs and recognizes the jwt processor in the pipeline.

### Assumptions

- The dnslookup processor is available and its structure is suitable for use as a skeleton.
- The collector builder supports custom processors and the jwt processor can be integrated without additional changes.
- Standard collector build and configuration practices are followed.

### Key Entities

- **JWT Processor**: Represents the new processor component, with attributes such as name, configuration, and enabled status.
- **Collector Build Config**: Represents the builder-config.yaml, listing enabled receivers, exporters, processors, and pipelines.

## Success Criteria *(mandatory)*

<!--
  ACTION REQUIRED: Define measurable success criteria.
  These must be technology-agnostic and measurable.
-->

### Measurable Outcomes

- **SC-001**: jwtprocessorスケルトンが作成され、コンパイルエラーなくビルドできる。
- **SC-002**: builder-config.yamlでjwtprocessorを有効化し、Goテストで `ocb --config builder-config.yaml` 実行・成果物の存在確認が自動化されている。
- **SC-003**: Goテストで `dist/otelcol-jwt --config otelcol-config.yaml` の実行・`sample_out.log` の内容検証（`signed.` の有無）が自動化されている。
- **SC-004**: .gitignoreで `dist/` と `test/sample_out.log` が除外されている。
