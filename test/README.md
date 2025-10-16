EdDSA ED25519 private key

```
openssl genpkey -algorithm ED25519 -out ed25519-private.pem
```

public key

```
openssl pkey -in ed25519-private.pem -pubout -out ed25519-public.pem
```

