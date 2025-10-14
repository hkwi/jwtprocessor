
# Data Model: opentelemetry-collector-processor

## JWT Processor Configuration

- **private_key_type**: string  
  - 許容値: "EdDSA_25519" または "RSA"
- **private_key**: string  
  - PEM形式の秘密鍵（EdDSA 25519 または RSA）

### 例

```yaml
jwtprocessor:
  private_key_type: "EdDSA_25519"
  private_key: |
    -----BEGIN PRIVATE KEY-----
    [base64 encoded key]
    -----END PRIVATE KEY-----
```

または

```yaml
jwtprocessor:
  private_key_type: "RSA"
  private_key: |
    -----BEGIN RSA PRIVATE KEY-----
    [base64 encoded key]
    -----END RSA PRIVATE KEY-----
```

## Validation Rules

- `private_key_type` は "EdDSA_25519" か "RSA" のいずれかでなければならない
- `private_key` は対応する形式のPEMであること

## State Transitions

- 設定: 未設定 → 設定済み → Collectorで利用可能
