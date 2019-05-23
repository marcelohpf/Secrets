package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
	"secrets/config"
	_ "secrets/server"
)

var rootCmd = &cobra.Command{
	Use:     "vault",
	Short:   "Vault is a simple secret managemnet tool.",
	Long:    "Vault is a simple secret management tool.",
	Version: "0.1.0",
	Run:     execute,
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&config.Verbose, "verbose", false, "Verbose mode")
	rootCmd.PersistentFlags().BoolVar(&config.Debug, "debug", false, "Debug mode")
	rootCmd.PersistentFlags().StringVar(&config.BoxPath, "box-path", config.DefaultBoxPath, "Box path")
	rootCmd.PersistentFlags().StringVar(&config.KeyPath, "key-path", config.DefaultKeyPath, "Cipher text path")
}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func execute(cmd *cobra.Command, args []string) {
	setupLog()
	log.Info("Running on [host]:[port]...", config.Verbose, config.Debug)
	log.Fatal(http.ListenAndServe(":8000", nil))
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
