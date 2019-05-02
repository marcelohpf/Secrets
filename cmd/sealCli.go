package cmd

import (
  "secrets/config"
  "secrets/crypto"
  "secrets/boxes"
  "github.com/spf13/cobra"
  log "github.com/sirupsen/logrus"
)

var boxCmd = &cobra.Command{
  Use: "box",
  Short: "Box",
  Long: "Manage your secret boxes.",
  Version: "0.0",
  Run: func(cmd *cobra.Command, args[]string){},
}

var encryptCmd = &cobra.Command{
  Use:   "seal",
  Short: "encrpyt",
  Long: "encrypt",
  Run: encrypt,
}
var decryptCmd = &cobra.Command{
  Use:   "unseal",
  Short: "decrpyt",
  Long: "decrypt",
  Run: decrypt,
}

func init(){
  boxCmd.PersistentFlags().StringVar(&config.ItemName, "item", "", "Item name to read/write content. Should be a single file on any type.")
  boxCmd.PersistentFlags().StringVar(&config.BoxName, "box", "", "Box name")
  boxCmd.PersistentFlags().StringVar(&config.BoxPath, "box-path", config.DefaultBoxPath, "Box path")
  boxCmd.PersistentFlags().StringVar(&config.KeyPath, "key-path", config.DefaultKeyPath, "Cipher text path")
  boxCmd.PersistentFlags().StringVar(&config.KeyName, "key", "", "Key name")
  encryptCmd.PersistentFlags().StringVar(&config.InFile, "in", "", "In file content to seal.")
  decryptCmd.PersistentFlags().StringVar(&config.OutFile, "out", "", "Out file content after unseal.")
  boxCmd.AddCommand(encryptCmd)
  boxCmd.AddCommand(decryptCmd)
  rootCmd.AddCommand(boxCmd)
}

func encrypt(cmd *cobra.Command, args[]string) {
  setupLog()

  key := crypto.GetKey(config.KeyPath, config.KeyName)
  // seal()
  text, err := boxes.ReadFromFile(config.InFile)
  if err != nil {
    log.Fatal(err.Error())
    panic("error")
  }
  ciphertext := crypto.Encrypt(text, key)
  switch config.BackendStorage {
    case "gdrive":
      err := boxes.GWriteBoxItem("", "", config.ItemName, ciphertext)
      if err != nil {
        log.Fatal(err.Error())
        panic("Item not retrieved from drive")
      }
    default:
      boxes.WriteBoxItem(config.BoxPath, config.BoxName, config.ItemName, ciphertext)
  }
}

func decrypt(cmd *cobra.Command, args[]string) {
  setupLog()
  key := crypto.GetKey(config.KeyPath, config.KeyName)
  // unseal
  var ciphertext string
  var err error
  switch config.BackendStorage {
    case "gdrive":
      ciphertext, err = boxes.GReadBoxItem("", "", config.ItemName)
    default:
      ciphertext, err = boxes.ReadBoxItem(config.BoxPath, config.BoxName, config.ItemName)
  }
  if err != nil {
    log.Fatal(err.Error())
    panic("Item not retrieved item")
  }
  text := crypto.Decrypt(ciphertext, key)
  boxes.WriteIntoFile(config.OutFile, text)
}

