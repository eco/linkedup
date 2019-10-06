package handler

import (
	"github.com/eco/longy/key-service/eventbrite"
	"github.com/eco/longy/key-service/mail"
	"github.com/eco/longy/key-service/masterkey"
	"github.com/eco/longy/key-service/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

// Router returns the root http Handler
func Router(eb *eventbrite.Session, mk *masterkey.Key, mc *mail.Client) http.Handler {
	r := mux.NewRouter()
	registerPing(r)
	registerRekey(r, eb, mk, mc)

	return middleware.LogHTTP(r)
}
