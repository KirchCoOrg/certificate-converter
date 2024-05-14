package main

import (
  "encoding/json"
  "crypto/rsa"
  "crypto/sha1"
  "crypto/x509"
  "encoding/base64"
  "encoding/pem"
  "fmt"
  "math/big"
  "github.com/docker/libtrust"
)

type JWK struct {
    Alg string `json:"alg"`
    Kty string `json:"kty"`
    Use string `json:"use"`
    X5c []string `json:"x5c"`
    Kid string `json:"kid"`
    N   string `json:"n,omitempty"`
    E   string `json:"e,omitempty"`
    X5t string `json:"x5t"`
}

type JWKS struct {
    Keys []JWK `json:"keys"`
}


func convertToJWKS(certData []byte) JWKS  {
  cert := safeParseX509Certificate(certData)

  rsaPublicKey, ok := cert.PublicKey.(*rsa.PublicKey)
  if !ok {
      fmt.Errorf("Public key is not of type RSA")
  }

  modulus := base64.RawURLEncoding.EncodeToString(rsaPublicKey.N.Bytes())
  exponent := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(rsaPublicKey.E)).Bytes())

  certSHA1 := sha1.Sum(cert.Raw)
  x5t := base64.RawURLEncoding.EncodeToString(certSHA1[:])

  x5c := base64.StdEncoding.EncodeToString(cert.Raw)


  publicKey, err := libtrust.FromCryptoPublicKey(rsaPublicKey)
  if err != nil {
      fmt.Errorf("Failed to create public key from RSA key: %v", err)
  }


  jwk := JWK{
      Alg: "RS256",
      Kty: "RSA",
      Use: "sig",
      X5c: []string{x5c},
      N:   modulus,
      E:   exponent,
      Kid: publicKey.KeyID(),
      X5t: x5t,
  }

  return JWKS{
    Keys: []JWK{jwk},
  }
}

func safeJwksToBytes(jwks JWKS) []byte {
  bytes, err := json.MarshalIndent(jwks, "", "  ")
  if err != nil {
    fmt.Errorf("Failed to marshal JWKS: %v", err)
  }
  return bytes
}

func safeParseX509Certificate(data []byte) *x509.Certificate {
  block, _ := pem.Decode(data)
  if block == nil {
      fmt.Errorf("Failed to decode PEM data")
  }

  cert, err := x509.ParseCertificate(block.Bytes)
  if err != nil {
      fmt.Errorf("Failed to parse certificate: %v", err)
  }
  return cert
}
