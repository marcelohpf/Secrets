package server

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
	"secrets/boxes"
	"secrets/config"
	"secrets/crypto"
)

type BoxArgs struct {
	Content  string        `json:"content"`
	BoxName  string        `json:"boxName"`
	ItemName string        `json:"itemName"`
	Key      string        `json:"key"`
	Token    *oauth2.Token `json:"token"`
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

func boxHandler(args BoxArgs) boxes.Vault {
	boxName := args.BoxName
	itemName := args.ItemName
	token := args.Token
	return boxes.Builder(config.BoxPath, boxName, itemName, token)
}

func sealHandler(w http.ResponseWriter, r *http.Request) {
	var args BoxArgs
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
	var args BoxArgs
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
	encoded, _ := json.Marshal(plainText)
	fmt.Fprintf(w, "{\"msg\": \"readed\", \"content\": %s}", string(encoded))
}
