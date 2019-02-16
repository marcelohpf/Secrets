package crypto

import (
  "crypto/rand"
  "log"
  "io"
  "encoding/base64"
  "os"
  "vault/config"
)

func GenerateKey() string {
  var logger = log.New(os.Stdout, "crypto ", log.Lshortfile)
  logger.Println("Validate Key Size")
  if config.SizeKey != 16 && config.SizeKey != 32 {
    logger.Fatal("Key should have size of 16 or 32 bytes.")
    panic("eerrr0r")
  }

  key := generate(config.SizeKey)
  logger.Println("Encoding generated key")
  return base64.StdEncoding.EncodeToString(key)
}

func generate(size int) []byte {
  var logger = log.New(os.Stdout, "crypto ", log.Lshortfile)
  logger.Println("Generating stream")
  data := make([]byte, size)
  if _, err := io.ReadFull(rand.Reader, data); err != nil {
    logger.Fatal(err.Error())
    panic(err.Error())
  }
  return data
}
