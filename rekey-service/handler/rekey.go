package handler

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/eco/longy/rekey-service/eventbrite"
	"github.com/eco/longy/rekey-service/mail"
	"github.com/eco/longy/rekey-service/masterkey"
	"github.com/eco/longy/util"
	"github.com/gorilla/mux"
	"net/http"
)

func registerRekey(r *mux.Router, eb eventbrite.Session, mk masterkey.Key, mc mail.Client) {
	r.HandleFunc("/rekey", rekey(eb, mk, mc)).Methods("GET")
}

// All core logic is implemented here. If there are plans to expand this service,
// logic (email retrieval, etc) can be lifted into http middleware to allow for better
// composability
func rekey(eb eventbrite.Session, mk masterkey.Key, mc mail.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type reqBody struct {
			AttendeeID int
			PublicKey  string
		}

		var body reqBody
		jsonDecoder := json.NewDecoder(r.Body)
		if err := jsonDecoder.Decode(&body); err != nil {
			errMsg := fmt.Sprintf("bad json request body: %s", err)
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}

		pubKeyBytes, err := hex.DecodeString(util.TrimHex(body.PublicKey))
		if err != nil {
			http.Error(w, "public key must be in hex format", http.StatusBadRequest)
			return
		}

		txBytes, err := mk.RekeyTransaction(body.AttendeeID, pubKeyBytes)
		if err != nil {
			http.Error(w, "internal error. try again", http.StatusInternalServerError)
			return
		}

		email, err := eb.EmailFromAttendeeID(body.AttendeeID)
		if err != nil {
			http.Error(w, "internal error. try again", http.StatusInternalServerError)
			return
		} else if len(email) == 0 {
			http.Error(w, "attendee id not present in the event", http.StatusNotFound)
			return
		}

		err = mc.SendRekeyEmail(email, txBytes)
		if err != nil {
			http.Error(w, "email error. try again", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
