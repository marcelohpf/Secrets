package cmd

import (
  "secrets/crypto"
  "secrets/config"
  "github.com/spf13/cobra"
)

var keyCmd = &cobra.Command{
  Use:   "keygen",
  Short: "Generate a key",
  Long: "Generate a key",
  Version: "0.0",
  Run: genKey,
}

func init(){
  keyCmd.PersistentFlags().IntVar(&config.SizeKey, "size-key", 32, "Define the size of generated key.")
  keyCmd.PersistentFlags().StringVar(&config.KeyPath, "key-path", config.DefaultKeyPath, "Cipher text path")
  keyCmd.PersistentFlags().StringVar(&config.KeyName, "key", "vault.key", "Key name")
  rootCmd.AddCommand(keyCmd)
}

func genKey(cmd *cobra.Command, args[]string) {
  setupLog()
  key := crypto.GenerateKey()
  crypto.SaveKey(config.KeyPath, config.KeyName, key)
}

