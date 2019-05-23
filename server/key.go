package server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"secrets/crypto"
)

func init() {
	genKeyRouter.registerRouter()
}

var genKeyRouter = Route{
	path:          "/key/generate",
	callback:      GenKeyHandler,
	documentation: "Generate a new key",
}

func GenKeyHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("Processing gen key http handler")
	key := crypto.GenerateKey()
	fmt.Fprintf(w, "{\"key\": \"%s\"}", key)
}
