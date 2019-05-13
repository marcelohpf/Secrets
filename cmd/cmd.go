package cmd

import (
  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
  "secrets/config"
)

var rootCmd = &cobra.Command{
  Use:   "vault",
  Short: "Vault is a simple secret managemnet tool.",
  Long: "Vault is a simple secret management tool.",
  Version: "0.1.0",
  Run: execute,
}

func init(){
  rootCmd.PersistentFlags().BoolVar(&config.Verbose, "verbose", false, "Verbose mode")
  rootCmd.PersistentFlags().BoolVar(&config.Debug, "debug", false, "Debug mode")
}
func Execute() {
  if err := rootCmd.Execute(); err != nil {
    log.Fatal(err)
  }
}

func execute(cmd *cobra.Command, args[]string){
  setupLog()
  log.Info("Running on [host]:[port]...", config.Verbose, config.Debug)
}

func setupLog() {
  if config.Verbose {
    log.SetLevel(log.InfoLevel)
  } else if config.Debug {
    log.SetLevel(log.TraceLevel)
  } else {
    log.SetLevel(log.WarnLevel)
  }
}
