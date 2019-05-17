package crypto

import (
  "crypto/rand"
  log "github.com/sirupsen/logrus"
  "io"
  "encoding/base64"
  "secrets/config"
)

func GenerateKey() string {
  log.Info("Generating a new key.")
  log.Debug("Validate Key Size")
  if config.SizeKey != 16 && config.SizeKey != 32 {
    log.Fatal("Key should have size of 16 or 32 bytes.")
    panic("Error")
  }

  key := generate(config.SizeKey)
  log.Debug("Encoding generated key")
  encodedKey := base64.StdEncoding.EncodeToString(key)
  log.Info("New key generated with success!")
  return encodedKey
}

func generate(size int) []byte {
  log.Info("Generating stream")
  data := make([]byte, size)
  if _, err := io.ReadFull(rand.Reader, data); err != nil {
    log.Fatal(err.Error())
    panic(err.Error())
  }
  return data
}
