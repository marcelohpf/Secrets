package cmd

import (
  "secrets/config"
  "secrets/crypto"
  "secrets/boxes"
  "github.com/spf13/cobra"
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
  Run: gencrypt,
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

func gencrypt(cmd *cobra.Command, args[]string) {
  setupLog()
  boxes.GReadBoxItem("", "", "secrets.vlt")
}

func encrypt(cmd *cobra.Command, args[]string) {
  setupLog()
  key := crypto.GetKey(config.KeyPath, config.KeyName)
  // seal()
  text := boxes.ReadFromFile(config.InFile)
  ciphertext := crypto.Encrypt(text, key)
  boxes.WriteBoxItem(config.BoxPath, config.BoxName, config.ItemName, ciphertext)
}

func decrypt(cmd *cobra.Command, args[]string) {
  setupLog()
  key := crypto.GetKey(config.KeyPath, config.KeyName)
  // unseal
  ciphertext := boxes.ReadBoxItem(config.BoxPath, config.BoxName, config.ItemName)
  text := crypto.Decrypt(ciphertext, key)
  boxes.WriteIntoFile(config.OutFile, text)
}

