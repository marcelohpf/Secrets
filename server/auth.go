package server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"secrets/boxes"
)

func init() {
	tokenRoute.registerRouter()
	authRoute.registerRouter()
}

var tokenRoute = Route{
	path:          "/goauth/token",
	callback:      tokenHandler,
	documentation: "route to receive the token from oauth",
}

var authRoute = Route{
	path:          "/goauth/auth",
	callback:      authHandler,
	documentation: "route to generate a google token",
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	authCode := query.Get("code")
	tokenJson, err := boxes.ExchangeToken(authCode)
	if err != nil {
		log.Error(err)
		fmt.Fprintf(w, "{\"msg\": \"problem to process the token\"}")
	}
	fmt.Fprintf(w, "{\"msg\": \"token received\", \"token\": \"%s\"}", tokenJson)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	authURL, err := boxes.GUrlAuth()
	if err == nil {
		fmt.Fprintf(w, "{\"redirect\": \"%s\"}", authURL)
	} else {
		log.Error(err)
		fmt.Fprintf(w, "{\"msg\": \"problem to redirect to gauth\"}")
	}
}
