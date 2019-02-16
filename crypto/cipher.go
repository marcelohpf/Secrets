package crypto

import (
  "crypto/aes"
  "crypto/cipher"
  "crypto/rand"
  "log"
  "io"
  "encoding/base64"
  "os"
)

func encode(text []byte) string {
  return base64.StdEncoding.EncodeToString(text)
}

func Encrypt(text, password string) string {
  var logger = log.New(os.Stdout, "crypto ", log.Lshortfile)
  logger.Println("initializing encryption")
  plaintext := []byte(text)
  key := decode(password)
  ciphertext, nonce := encrypt(plaintext, key)
  cipher := append(nonce, ciphertext...)
  return encode(cipher)

}

func encrypt(text, key []byte) ([]byte, []byte) {

  var logger = log.New(os.Stdout, "crypto ", log.Lshortfile)

  logger.Println("initializing key")

  logger.Println("texto to bytes")
  plaintext := []byte(text)
  logger.Println("initialize block")
  block, err := aes.NewCipher(key)

  if err != nil {
    logger.Fatal(err.Error())
    panic(err.Error())
  }

  nonce := make([]byte, 12)
  if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
    logger.Fatal(err.Error())
    panic(err.Error())
  }

  aesgcm, err := cipher.NewGCM(block)
  if err != nil {
    logger.Fatal(err.Error())
    panic(err.Error())
  }

  ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

  return ciphertext, nonce
}
