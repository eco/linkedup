package handler

import (
	"github.com/eco/longy/rekey-service/eventbrite"
	"github.com/eco/longy/rekey-service/mail"
	"github.com/eco/longy/rekey-service/masterkey"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func registerRekey(r *mux.Router, eb eventbrite.Session, mk masterkey.Key, mc mail.Client) {
	r.HandleFunc("/rekey/{id:[0-9]+}", rekey(eb, mk, mc)).Methods("GET")
}

// All core logic is implemented here. If there are plans to expand this service,
// logic (email retrieval, etc) can be lifted into http middleware to allow for better
// composability
func rekey(eb eventbrite.Session, mk masterkey.Key, mc mail.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nonceStr := r.FormValue("nonce")
		if len(nonceStr) == 0 {
			http.Error(w, "nonce expected as a url parameter", http.StatusBadRequest)
			return
		}

		// retrieve id & nonce from the url path
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "attendee id must be in decimal format", http.StatusBadRequest)
			return
		}
		nonce, err := strconv.Atoi(nonceStr)
		if err != nil {
			http.Error(w, "nonce must be in decimal format", http.StatusBadRequest)
			return
		}

		sig, err := mk.RekeySignature(id, nonce)
		if err != nil {
			http.Error(w, "internal error. try again", http.StatusInternalServerError)
			return
		}

		email, err := eb.EmailFromAttendeeID(id)
		if err != nil {
			http.Error(w, "internal error. try again", http.StatusInternalServerError)
			return
		} else if len(email) == 0 {
			http.Error(w, "attendee id not present in the event", http.StatusNotFound)
			return
		}

		err = mc.SendRekeyEmail(email, sig)
		if err != nil {
			http.Error(w, "email error. try again", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
