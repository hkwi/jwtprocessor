EdDSA ED25519 private key

```
openssl genpkey -algorithm ED25519 -out ed25519-private.pem
```

public key

```
openssl pkey -in ed25519-private.pem -pubout -out ed25519-public.pem
```


# RSA

```
openssl genpkey -algorithm RSA -out rsa_key_pkcs1.pem -pkeyopt rsa_keygen_bits:2048
```

```
openssl rsa -in rsa_key_pkcs1.pem -pubout
```

