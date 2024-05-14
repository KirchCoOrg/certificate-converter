package main

import (
  "encoding/base64"
  "encoding/json"
  "fmt"
  "log"
)

type AcmeData struct {
  LetsEncrypt struct {
    Certificates []struct {
      Domain struct {
        Main string `json:"main"`
      } `json:"domain"`
      Certificate string `json:"certificate"`
      Key         string `json:"key"`
      Store       string `json:"Store"`
    } `json:"Certificates"`
  } `json:"LetsEncrypt"`
}

type Keypair struct {
  Certificate []byte
  Key         []byte
}

func getKeypair(acme *AcmeData, domain string) (Keypair, error) {
  for _, cert := range (*acme).LetsEncrypt.Certificates { // Iterate over all certificates
    if cert.Domain.Main == domain {

      // Decode the base64 encoded certificate and key values
      certificateDecoded, certErr := base64.StdEncoding.DecodeString(cert.Certificate)
      if certErr != nil {
        return Keypair{}, fmt.Errorf("Failed to decode certificate: %v", certErr)
      }
      keyDecoded, keyErr := base64.StdEncoding.DecodeString(cert.Key)
      if keyErr != nil {
        return Keypair{}, fmt.Errorf("Failed to decode key: %v", keyErr)
      }

      return Keypair{Certificate: certificateDecoded, Key: keyDecoded}, nil
    }
  }
  return Keypair{}, fmt.Errorf("No keypair found for domain: " + domain)
}

func safeGetKeypair(data *AcmeData, domain string) Keypair {
  keypair, err := getKeypair(data, domain)
  if err != nil {
    log.Fatalf("Failed to get keypair for domain: %v", err)
  }
  return keypair
}

func safeParseAcme(data []byte) AcmeData {
  var output AcmeData
  if err := json.Unmarshal(data, &output); err != nil {
    log.Fatalf("Failed to unmarshal JSON: %v", err)
  }
  return output
}