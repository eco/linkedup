package handler

import (
	"encoding/json"
	"fmt"
	"github.com/eco/longy/key-service/types"
	"github.com/eco/longy/util"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"net/http"
)

func registerRekey(r *mux.Router, svc types.Service) {
	r.HandleFunc("/rekey", rekey(svc)).Methods("GET")
}

// All core logic is implemented here. If there are plans to expand this service,
// logic (email retrieval, etc) can be lifted into http middleware to allow for better
// composability
func rekey(svc types.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type reqBody struct {
			AttendeeID string
			PublicKey  secp256k1.PubKeySecp256k1
		}

		var body reqBody
		jsonDecoder := json.NewDecoder(r.Body)
		if err := jsonDecoder.Decode(&body); err != nil {
			errMsg := fmt.Sprintf("bad json request body: %s", err)
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}

		secret, commitment := util.CreateCommitment()
		err := svc.MasterKey().SendRekeyTransaction(body.AttendeeID, body.PublicKey, commitment)
		if err != nil {
			http.Error(w, "internal error. try again", http.StatusInternalServerError)
			return
		}

		email, err := svc.Eventbrite().EmailFromAttendeeID(body.AttendeeID)
		if err != nil {
			http.Error(w, "internal error. try again", http.StatusInternalServerError)
			return
		} else if len(email) == 0 {
			http.Error(w, "attendee id not present in the event", http.StatusNotFound)
			return
		}

		err = svc.MailClient().SendRekeyEmail(email, secret)
		if err != nil {
			http.Error(w, "email error. try again", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
