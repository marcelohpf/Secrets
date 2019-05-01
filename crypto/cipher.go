package crypto

import (
  "crypto/aes"
  "crypto/cipher"
  "crypto/rand"
  log "github.com/sirupsen/logrus"
  "io"
  "encoding/base64"
)

func encode(text []byte) string {
  return base64.StdEncoding.EncodeToString(text)
}

func Encrypt(text, password string) string {
  log.Info("initializing encryption.")
  plaintext := []byte(text)
  key := decode(password)
  ciphertext, nonce := encrypt(plaintext, key)
  cipher := append(nonce, ciphertext...)
  return encode(cipher)

}

func encrypt(text, key []byte) ([]byte, []byte) {
  log.Info("initializing key.")
  log.Info("text to bytes.")
  plaintext := []byte(text)
  log.Info("initialize block.")
  block, err := aes.NewCipher(key)

  if err != nil {
    log.Fatal(err.Error())
    panic(err.Error())
  }

  nonce := make([]byte, 12)
  if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
    log.Fatal(err.Error())
    panic(err.Error())
  }

  aesgcm, err := cipher.NewGCM(block)
  if err != nil {
    log.Fatal(err.Error())
    panic(err.Error())
  }

  ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

  return ciphertext, nonce
}
