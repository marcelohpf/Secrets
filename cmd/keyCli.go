package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
	"os/exec"
	"runtime"
	"secrets/boxes"
	"secrets/config"
	"secrets/crypto"
)

var keyCmd = &cobra.Command{
	Use:   "genkey",
	Short: "Generate a key",
	Long:  "Generate a key",
	Run:   genKey,
}

var gtokenCmd = &cobra.Command{
	Use:   "gdrive-token-refresh",
	Short: "Refresh drive token",
	Long:  "Regenerate the auth token for google drive api",
	Run:   fetchToken,
}

func init() {
	keyCmd.PersistentFlags().IntVarP(&config.SizeKey, "size-key", "s", 32, "Define the size of generated key.")
	keyCmd.PersistentFlags().StringVarP(&config.KeyName, "key-name", "n", "secrets", "Key name")
	rootCmd.AddCommand(keyCmd)
	rootCmd.AddCommand(gtokenCmd)
}

func fetchToken(cmd *cobra.Command, args []string) {
	setupLog()
	r, err := http.Get("http://" + config.Server + ":" + config.Port + "/goauth/auth")
	if err != nil {
		log.Fatal(err)
		return
	}
	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Fatal(err)
		return
	}

	err = openUrl(data["redirect"])

	var token string
	fmt.Printf("Enter here the token from web browser")
	if _, err := fmt.Scan(&token); err != nil {
		log.Fatal(err)
	}
	if err := boxes.WriteIntoFile(config.TokenFile, string(token)); err != nil {
		log.Fatal(err)
	}
}

func genKey(cmd *cobra.Command, args []string) {
	setupLog()

	r, err := http.Get("http://" + config.Server + ":" + config.Port + "/key/generate")
	if err != nil {
		log.Fatal(err)
		return
	}

	var key map[string]string
	if err = json.NewDecoder(r.Body).Decode(&key); err != nil {
		log.Fatal(err)
		return
	}
	log.Debug(key)

	log.Info(config.KeyName)
	if err := crypto.SaveKey(config.KeyPath, config.KeyName, key["key"]); err != nil {
		log.Fatal(err)
	}
}

func openUrl(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = errors.New("unsupported platform")
	}
	return err
}
