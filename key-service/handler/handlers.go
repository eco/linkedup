package handler

import (
	"github.com/eco/longy/key-service/eventbrite"
	"github.com/eco/longy/key-service/mail"
	"github.com/eco/longy/key-service/masterkey"
	"github.com/eco/longy/key-service/middleware"
	"github.com/eco/longy/key-service/models"
	"github.com/eco/longy/x/longy/client/rest"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

var log = logrus.WithField("module", "handler")

// Router returns the root http Handler
func Router(
	eb *eventbrite.Session,
	mk *masterkey.MasterKey,
	db *models.DatabaseContext,
	mc mail.Client) http.Handler {

	r := mux.NewRouter()
	// IMPORTANT: you must specify an OPTIONS method matcher for the middleware to set CORS headers
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(rest.CorsMiddleware)

	registerPing(r)
	registerKey(r, eb, mk, db, mc)
	registerEmailManual(r, db, eb, mc)
	registerInfo(r, db, mc)
	registerIDToAddress(r)

	return middleware.LogHTTP(r)
}
