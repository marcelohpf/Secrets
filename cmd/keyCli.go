package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"secrets/boxes"
	"secrets/config"
	"secrets/crypto"
)

var keyCmd = &cobra.Command{
	Use:     "keygen",
	Short:   "Generate a key",
	Long:    "Generate a key",
	Version: "0.1",
	Run:     genKey,
}

var gtokenCmd = &cobra.Command{
	Use:     "gdrive-token-refresh",
	Short:   "Refresh drive token",
	Long:    "Regenerate the auth token for google drive api",
	Version: "0.0",
	Run: func(cmd *cobra.Command, args []string) {
		setupLog()
		_, err := boxes.GRefreshAuth()
		if err != nil {
			log.Fatal(err.Error())
			panic("Could not refresh token for google drive api")
		}
	},
}

func init() {
	keyCmd.PersistentFlags().IntVar(&config.SizeKey, "size-key", 32, "Define the size of generated key.")
	keyCmd.PersistentFlags().StringVar(&config.KeyName, "key", "vault.key", "Key name")
	rootCmd.AddCommand(keyCmd)
	rootCmd.AddCommand(gtokenCmd)
}

func genKey(cmd *cobra.Command, args []string) {
	setupLog()
	key := crypto.GenerateKey()
	if err := crypto.SaveKey(config.KeyPath, config.KeyName, key); err != nil {
		log.Fatal(err)
	}
}
