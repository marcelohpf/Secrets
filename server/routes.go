package server

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Route struct {
	path          string `json:"path"`
	documentation string `json:"documentation"`
	callback      func(w http.ResponseWriter, r *http.Request)
}

func (r Route) registerRouter() {
	log.Debug("Registering %s", r.path)
	http.HandleFunc(r.path, r.callback)
}
