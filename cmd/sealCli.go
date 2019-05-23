package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"secrets/boxes"
	"secrets/config"
	"secrets/crypto"
)

var boxCmd = &cobra.Command{
	Use:     "box",
	Short:   "Box",
	Long:    "Manage your secret boxes.",
	Version: "0.1",
	Run:     func(cmd *cobra.Command, args []string) {},
}

var encryptCmd = &cobra.Command{
	Use:   "seal",
	Short: "encrpyt",
	Long:  "encrypt",
	Run:   encrypt,
}
var decryptCmd = &cobra.Command{
	Use:   "unseal",
	Short: "decrpyt",
	Long:  "decrypt",
	Run:   decrypt,
}

func init() {
	boxCmd.PersistentFlags().StringVar(&config.ItemName, "item", "", "Item name to read/write content. Should be a single file on any type.")
	boxCmd.PersistentFlags().StringVar(&config.BoxName, "box", "", "Box name")
	boxCmd.PersistentFlags().StringVar(&config.KeyName, "key", "", "Key name")
	encryptCmd.PersistentFlags().StringVar(&config.InFile, "in", "", "In file content to seal.")
	decryptCmd.PersistentFlags().StringVar(&config.OutFile, "out", "", "Out file content after unseal.")
	boxCmd.AddCommand(encryptCmd)
	boxCmd.AddCommand(decryptCmd)
	rootCmd.AddCommand(boxCmd)
}

func encrypt(cmd *cobra.Command, args []string) {
	setupLog()

	key, err := crypto.GetKey(config.KeyPath, config.KeyName)
	if err != nil {
		log.Fatal(err.Error())
		panic("error")
	}

	// seal
	text, err := boxes.ReadFromFile(config.InFile)

	if err != nil {
		log.Fatal(err.Error())
		panic("error")
	}

	ciphertext := crypto.Encrypt(text, key)

	box := boxes.Builder(config.BoxPath, config.BoxName, config.ItemName)
	err = box.WriteBoxItem(ciphertext)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func decrypt(cmd *cobra.Command, args []string) {
	setupLog()

	key, err := crypto.GetKey(config.KeyPath, config.KeyName)

	if err != nil {
		log.Fatal(err.Error())
		panic("error")
	}

	box := boxes.Builder(config.BoxPath, config.BoxName, config.ItemName)

	// unseal
	ciphertext, err := box.ReadBoxItem()

	if err != nil {
		log.Fatal(err.Error())
		panic("Item not retrieved item")
	}

	text := crypto.Decrypt(ciphertext, key)
	boxes.WriteIntoFile(config.OutFile, text)
}
