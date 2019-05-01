package crypto

import (
  "crypto/aes"
  "crypto/cipher"
  log "github.com/sirupsen/logrus"
  "encoding/base64"
)

func Decrypt(ciphertext, password string) (string) {

  log.Info("Starting the decrypt.")
  byte_text := decode(ciphertext)
  log.Info("Extract nonce.")
  nonce := byte_text[:12]
  text := byte_text[12:]
  key := decode(password)
  log.Info("Decrypt plain text.")
  plaintext := decrypt(text, nonce, key)
  return string(plaintext)
}

func decode(text string) []byte {
  decoded_text, err := base64.StdEncoding.DecodeString(text)
  if err != nil {
    panic(err.Error())
  }
  return decoded_text
}


func decrypt(ciphertext, nonce, key []byte) []byte {
  log.Info("Initialize block")
  block, err := aes.NewCipher(key)

  if err != nil {
    log.Fatal(err.Error())
    panic(err.Error())
  }

  log.Info("Initialize cipher")
  aesgcm, err := cipher.NewGCM(block)

  if err != nil {
    log.Fatal(err.Error())
    panic(err.Error())
  }

  log.Info("Open gcm")
  plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
  if err != nil {
    log.Fatal(err.Error())
    panic(err.Error())
  }

  return plaintext
}
