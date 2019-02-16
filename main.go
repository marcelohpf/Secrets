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
    key := crypto.GetKey()
    text := boxes.ReadFileBox(config.TextPath)
    ciphertext := crypto.Encrypt(text, key)
    boxes.WriteFileBox(config.CiphertextPath, ciphertext)
  } else if config.Decrypt {
    key := crypto.GetKey()
    text := boxes.ReadFileBox(config.CiphertextPath)
    ciphertext := crypto.Decrypt(text, key)
    fmt.Println(ciphertext)
  } else if config.GenerateKey {
    key := crypto.GenerateKey()
    crypto.SaveKey(key)
  }
}
