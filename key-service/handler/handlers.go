package handler

import (
	"github.com/eco/longy/key-service/eventbrite"
	"github.com/eco/longy/key-service/mail"
	"github.com/eco/longy/key-service/masterkey"
	"github.com/eco/longy/key-service/middleware"
	"github.com/eco/longy/key-service/models"
	"github.com/gorilla/mux"
	"net/http"
)

// Router returns the root http Handler
func Router(
	eb *eventbrite.Session,
	mk *masterkey.MasterKey,
	db *models.DatabaseContext,
	mc mail.Client) http.Handler {

	r := mux.NewRouter()
	registerPing(r)
	registerKey(r, eb, mk, db, mc)

	return middleware.LogHTTP(r)
}
