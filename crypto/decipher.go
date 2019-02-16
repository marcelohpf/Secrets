package crypto

import (
  "crypto/aes"
  "crypto/cipher"
  "log"
  "encoding/base64"
  "os"
)

func decode(text string) []byte {
  decoded_text, err := base64.StdEncoding.DecodeString(text)
  if err != nil {
    panic(err.Error())
  }
  return decoded_text
}

func Decrypt(ciphertext, password string) (string) {

  var logger = log.New(os.Stdout, "crypto ", log.Lshortfile)

  logger.Println("Starting the decrypt")
  byte_text := decode(ciphertext)
  logger.Println("Extract nonce")
  nonce := byte_text[:12]
  text := byte_text[12:]
  key := decode(password)
  logger.Println("Decrypt plain text")
  plaintext := decrypt(text, nonce, key)
  return string(plaintext)
}

func decrypt(ciphertext, nonce, key []byte) []byte {
  var logger = log.New(os.Stdout, "crypto ", log.Lshortfile)


  logger.Println("Initialize block")
  block, err := aes.NewCipher(key)

  if err != nil {
    logger.Fatal(err.Error())
    panic(err.Error())
  }

  logger.Println("Initialize cipher")
  aesgcm, err := cipher.NewGCM(block)

  if err != nil {
    logger.Fatal(err.Error())
    panic(err.Error())
  }

  logger.Println("Open gcm")
  plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
  if err != nil {
    logger.Fatal(err.Error())
    panic(err.Error())
  }

  return plaintext
}
