package handler

import (
	"encoding/json"
	"fmt"
	"github.com/eco/longy/key-service/eventbrite"
	"github.com/eco/longy/key-service/mail"
	"github.com/eco/longy/key-service/masterkey"
	"github.com/eco/longy/util"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"net/http"
)

func registerKey(r *mux.Router, eb *eventbrite.Session, mk *masterkey.MasterKey, mc *mail.Client) {
	r.HandleFunc("/key", key(eb, mk, mc)).Methods("POST")
}

// All core logic is implemented here. If there are plans to expand this service,
// logic (email retrieval, etc) can be lifted into http middleware to allow for better
// composability
func key(eb *eventbrite.Session, mk *masterkey.MasterKey, mc *mail.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type reqBody struct {
			AttendeeID string                    `json:"attendee_id"`
			PubKey     secp256k1.PubKeySecp256k1 `json:"pubkey"`
		}

		var body reqBody
		jsonDecoder := json.NewDecoder(r.Body)
		if err := jsonDecoder.Decode(&body); err != nil {
			errMsg := fmt.Sprintf("bad json request body: %s", err)
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}

		secret, commitment := util.CreateCommitment()
		err := mk.SendKeyTransaction(body.AttendeeID, body.PubKey, commitment)
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

		err = mc.SendRedirectEmail(email, secret)
		if err != nil {
			http.Error(w, "email error. try again", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
