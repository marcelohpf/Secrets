package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
	"secrets/boxes"
	"secrets/config"
	"secrets/crypto"
	"secrets/server"
)

var boxCmd = &cobra.Command{
	Use:     "box",
	Short:   "Box",
	Long:    "Manage your secret boxes.",
	Version: "0.1",
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
	boxCmd.PersistentFlags().StringVar(&config.BoxKeyName, "key", "", "Key name")
	encryptCmd.PersistentFlags().StringVar(&config.InFile, "in", "", "In file content to seal.")
	decryptCmd.PersistentFlags().StringVar(&config.OutFile, "out", "", "Out file content after unseal.")
	boxCmd.AddCommand(encryptCmd)
	boxCmd.AddCommand(decryptCmd)
	rootCmd.AddCommand(boxCmd)
}

func encrypt(cmd *cobra.Command, args []string) {
	setupLog()

	key, err := crypto.GetKey(config.KeyPath, config.BoxKeyName)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	text, err := boxes.ReadFromFile(config.InFile)

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	var body server.BoxArgs

	switch config.BackendStorage {
	case "gdrive":
		token, err := boxes.TokenFromFile(config.TokenFile)
		if err != nil {
			log.Fatal(err)
			return
		}
		body = server.BoxArgs{Key: key, Content: text, BoxName: config.BoxName, ItemName: config.ItemName, Token: token}
	default:
		body = server.BoxArgs{Key: key, Content: text, BoxName: config.BoxName, ItemName: config.ItemName}
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
		return
	}

	r, err := http.Post("http://"+config.Server+":"+config.Port+"/box/seal", "application/json", bytes.NewBuffer(jsonBody))

	if err != nil {
		log.Fatal(err)
		return
	}

	var data map[string]string
	if err = json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Fatal(err)
		return
	}
	log.Info(data)
}

func decrypt(cmd *cobra.Command, args []string) {
	setupLog()

	key, err := crypto.GetKey(config.KeyPath, config.BoxKeyName)

	if err != nil {
		log.Fatal(err.Error())
		panic("error")
	}

	var body server.BoxArgs
	switch config.BackendStorage {
	case "gdrive":
		token, err := boxes.TokenFromFile(config.TokenFile)
		if err != nil {
			log.Fatal(err)
			return
		}
		body = server.BoxArgs{Key: key, BoxName: config.BoxName, ItemName: config.ItemName, Token: token}
	default:
		body = server.BoxArgs{Key: key, BoxName: config.BoxName, ItemName: config.ItemName}
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
		return
	}

	r, err := http.Post("http://"+config.Server+":"+config.Port+"/box/unseal", "application/json", bytes.NewBuffer(jsonBody))

	if err != nil {
		log.Fatal(err)
		panic("Item not retrieved item")
	}

	var data map[string]string

	if err = json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Fatal(err)
		return
	}

	if config.OutFile != "" {
		boxes.WriteIntoFile(config.OutFile, data["content"])
	} else {
		fmt.Println(data["content"])
	}
}
