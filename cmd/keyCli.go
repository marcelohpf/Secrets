package cmd

import (
  "secrets/crypto"
  "secrets/boxes"
  "secrets/config"
  "github.com/spf13/cobra"
  log "github.com/sirupsen/logrus"
)

var keyCmd = &cobra.Command{
  Use:   "keygen",
  Short: "Generate a key",
  Long: "Generate a key",
  Version: "0.0",
  Run: genKey,
}

var gtokenCmd = &cobra.Command{
  Use:   "gdrive-token-refresh",
  Short: "Refresh drive token",
  Long: "Regenerate the auth token for google drive api",
  Version: "0.0",
  Run: func(cmd *cobra.Command, args[]string) {
    setupLog()
    err := boxes.GRefreshAuth()
    if err != nil {
      log.Fatal(err.Error())
      panic("Could not refresh token for google drive api")
    }
  },
}

func init(){
  keyCmd.PersistentFlags().IntVar(&config.SizeKey, "size-key", 32, "Define the size of generated key.")
  keyCmd.PersistentFlags().StringVar(&config.KeyPath, "key-path", config.DefaultKeyPath, "Cipher text path")
  keyCmd.PersistentFlags().StringVar(&config.KeyName, "key", "vault.key", "Key name")
  rootCmd.AddCommand(keyCmd)
  rootCmd.AddCommand(gtokenCmd)
}

func genKey(cmd *cobra.Command, args[]string) {
  setupLog()
  key := crypto.GenerateKey()
  crypto.SaveKey(config.KeyPath, config.KeyName, key)
}

