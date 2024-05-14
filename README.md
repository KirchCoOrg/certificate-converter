# Certificate Converter

This is a simple tool to convert certificates from Traefik to X509 Certificate and Key, and JSON Web Key Sets.

Te project is currently an MVP and shall be generalized to allow converting certificates in multiple ways.


## Usage

```yaml
  certificate-converter:
    image: kirchco/certificate-converter
    environment:
      CERT_DOMAIN: "{{ hostnames[0] }}"
    volumes:
      - traefik-letsencrypt:/letsencrypt:ro
      - token-certificate:/cert
```

### Environment Variables
- `CERT_DOMAIN`: The domain of the TLS certificate to convert.

### Volumes
- `/letsencrypt`: The directory containing the Traefik certificates (`acme.json`)
- `/cert`: The directory to write the converted certificates to (`<domain>.pem`, `<domain>.key`, `jwks.json`)