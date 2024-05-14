package main

import (
  "log"
  "os"
  "time"
  "github.com/fsnotify/fsnotify"
)

var (
  certDomain = safeGetenv("CERT_DOMAIN")
  acmeFilePath = "/letsencrypt/acme.json"
  watcher *fsnotify.Watcher
)

func main() {
  // Check if acme file exists, If  not, wait one minute and retry
  for !safeStat(acmeFilePath) {
    time.Sleep(1 * time.Minute)
  }

  // Perform initial read and extract of acme.json
  extractAndSave()

  // Create watcher for the acme file
  watcher = safeNewWatcher()
  watcher.Add(acmeFilePath)
  defer watcher.Close()

  // Run deduplicated watcher for write events on acme file
  go deduplicatedWriteListener(extractAndSave)

  // Block forever
  <-make(chan struct{})
}

func extractAndSave() {
  // read and parse acme.json
  data := safeParseAcme(safeReadFile(acmeFilePath))

  // find domain and export resprective keypair (certificate and key)
  targetKeypair := safeGetKeypair(&data, certDomain)

  safeWriteFile(targetKeypair.Certificate, "/cert/" + certDomain + ".pem")
  safeWriteFile(targetKeypair.Key, "/cert/" + certDomain + ".key")

  jwks := convertToJWKS(targetKeypair.Certificate)
  safeWriteFile(safeJwksToBytes(jwks), "/cert/" + certDomain + ".jwks")
}



func safeGetenv(key string) string {
  value := os.Getenv(key)
  if value == "" {
    log.Fatalf("Environment variable %s is not set.", key)
  }
  return value
}

func safeStat(path string) bool { // Check if file in <path> exists
  _, err := os.Stat(path)
  if os.IsNotExist(err) {
    log.Printf("File not found at %s", path)
  }
  return !os.IsNotExist(err)
}

func safeReadFile(filename string) []byte {
  content, err := os.ReadFile(filename)
  if err != nil {
    log.Fatalf("Failed to read file: %s %v", filename, err)
  }
  return content
}

func safeWriteFile(data []byte, path string) {
  if err := os.WriteFile(path, data, 0644); err != nil {
    log.Fatalf("Failed to write file: %s %v", path, err)
  } else {
    log.Printf("Wrote file: %s", path)
  }
}

func safeNewWatcher() *fsnotify.Watcher {
  w, err := fsnotify.NewWatcher()
  if err != nil {
    log.Fatalf("Failed to create watcher: %v", err)
  }
  return w
}




