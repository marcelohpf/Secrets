package cmd

import (
  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
  Use:   "vault",
  Short: "Vault is a simple secret managemnet tool.",
  Long: "Vault is a simple secret management tool.",
  Version: "0.1.0",
  Run: execute,
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    log.Fatal(err)
  }
}

func execute(cmd *cobra.Command, args[]string){
  fmt.Println("Use --help")
  log.Info("Running on [host]:[port]...")
}
