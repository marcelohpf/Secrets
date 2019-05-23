package server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func init() {
	homeRouter.registerRouter()
}

var homeRouter = Route{
	path:          "/",
	callback:      homeHandler,
	documentation: "ping here",
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("Server check")
	fmt.Fprintf(w, "{\"msg\": \"ok\"}")
}
