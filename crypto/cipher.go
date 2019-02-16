package crypto

import (
  "crypto/aes"
  "crypto/cipher"
  "crypto/rand"
  "encoding/hex"
  "log"
  "io"
  "encoding/base64"
  "os"
)

func Encode(text []byte) string {
  return base64.StdEncoding.EncodeToString(text)
}

func Encrypt(text string) string {
  var logger = log.New(os.Stdout, "crypto ", log.Lshortfile)
  logger.Println("initializing encryption")
  plaintext := []byte(text)
  ciphertext, nonce := encrypt(plaintext)
  cipher := append(nonce, ciphertext...)
  return Encode(cipher)

}

func encrypt(text []byte) ([]byte, []byte) {

  var logger = log.New(os.Stdout, "crypto ", log.Lshortfile)

  logger.Println("initializing key")
  key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")

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
