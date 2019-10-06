package handler

import (
	"fmt"
	"github.com/eco/longy/key-service/eventbrite"
	"github.com/eco/longy/key-service/mail"
	"github.com/eco/longy/key-service/masterkey"
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy"
	"github.com/gorilla/mux"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	"io/ioutil"
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
			AttendeeID string          `json:"attendee_id"`
			PubKey     tmcrypto.PubKey `json:"pubkey"`
		}
		bz, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "could not read request body", http.StatusBadRequest)
			return
		}

		cdc := longy.ModuleCdc
		var body reqBody
		if err := cdc.UnmarshalJSON(bz, &body); err != nil {
			errMsg := fmt.Sprintf("bad json request body: %s", err)
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}

		secret, commitment := util.CreateCommitment()
		if err := mk.SendKeyTransaction(body.AttendeeID, body.PubKey, commitment); err != nil {
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
