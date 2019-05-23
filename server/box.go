package server

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"secrets/boxes"
	"secrets/config"
	"secrets/crypto"
)

type boxArgs struct {
	Content  string `json:"content"`
	BoxName  string `json:"boxName"`
	ItemName string `json:"itemName"`
	Key      string `json:"key"`
}

func init() {
	sealRoute.registerRouter()
	unsealRoute.registerRouter()
}

var sealRoute = Route{
	path:          "/box/seal",
	callback:      sealHandler,
	documentation: "seal a content into a box",
}

var unsealRoute = Route{
	path:          "/box/unseal",
	callback:      unsealHandler,
	documentation: "seal a content into a box",
}

func boxHandler(args boxArgs) boxes.Vault {
	boxName := args.BoxName
	itemName := args.ItemName
	return boxes.Builder(config.BoxPath, boxName, itemName)
}

func sealHandler(w http.ResponseWriter, r *http.Request) {
	var args boxArgs
	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		log.Debug("Problem decoding body request")
	}
	log.Debug(args)

	ciphertext := crypto.Encrypt(args.Content, args.Key)

	box := boxHandler(args)
	if err := box.WriteBoxItem(ciphertext); err != nil {
		log.Fatal(err.Error())
	}

	fmt.Fprintf(w, "{\"msg\": \"saved\"}")
}

func unsealHandler(w http.ResponseWriter, r *http.Request) {
	var args boxArgs
	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		log.Debug("Problem decoding body request")
	}

	log.Debug(args.Key)

	box := boxHandler(args)
	content, err := box.ReadBoxItem()

	if err != nil {
		log.Fatal(err.Error())
	}

	plainText := crypto.Decrypt(content, args.Key)
	fmt.Fprintf(w, "{\"msg\": \"readed\", \"content\": \"%s\"}", plainText)
}
