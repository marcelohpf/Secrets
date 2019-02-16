package main

import (
  "fmt"
  "vault/config"
  "vault/crypto"
  "vault/boxes"
)

func main(){
  fmt.Println("Personal Vault Service")
  execute()
}

func execute() {
  if config.Encrypt {
    text := boxes.ReadFileBox(config.TextPath)
    ciphertext := crypto.Encrypt(text)
    boxes.WriteFileBox(config.CiphertextPath, ciphertext)
  } else if config.Decrypt {
    text := boxes.ReadFileBox(config.CiphertextPath)
    ciphertext := crypto.Decrypt(text)
    fmt.Println(ciphertext)
  }
}
