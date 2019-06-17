package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"secrets/config"
)

var rootCmd = &cobra.Command{
	Use:     "vault",
	Short:   "Vault is a simple secret managemnet tool.",
	Long:    "Vault is a simple secret management tool.",
	Version: "0.2.0",
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&config.Verbose, "verbose", "v", false, "Verbose mode")
	rootCmd.PersistentFlags().BoolVarP(&config.Debug, "debug", "", false, "Debug mode")
	rootCmd.PersistentFlags().StringVarP(&config.KeyPath, "key-path", "", config.DefaultKeyPath, "Cipher text path")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
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
