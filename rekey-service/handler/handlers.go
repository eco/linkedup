package handler

import (
	"github.com/eco/longy/rekey-service/eventbrite"
	"github.com/eco/longy/rekey-service/mail"
	"github.com/eco/longy/rekey-service/masterkey"
	"github.com/eco/longy/rekey-service/middleware"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

var log = logrus.WithField("module", "router")

// Router returns the root http Handler
func Router(eb eventbrite.Session, mk masterkey.Key, mc mail.Client) http.Handler {
	r := mux.NewRouter()
	registerPing(r)
	registerRekey(r, eb, mk, mc)

	return middleware.LogHTTP(r)
}
