package handler

import (
	"encoding/json"
	"fmt"
	"github.com/eco/longy/key-service/eventbrite"
	"github.com/eco/longy/key-service/mail"
	"github.com/eco/longy/key-service/masterkey"
	"github.com/eco/longy/key-service/models"
	"github.com/eco/longy/util"
	"github.com/gorilla/mux"
	"net/http"
)

func registerKey(
	r *mux.Router,
	eb *eventbrite.Session,
	mk *masterkey.MasterKey,
	db *models.DatabaseContext,
	mc *mail.Client) {

	r.HandleFunc("/key", key(eb, mk, db, mc)).Methods("POST")
	r.HandleFunc("/key/{email}", keyGetter(db)).Methods("GET")
}

// All core logic is implemented here. If there are plans to expand this service,
// logic (email retrieval, etc) can be lifted into http middleware to allow for better
// composability
func key(eb *eventbrite.Session, mk *masterkey.MasterKey, db *models.DatabaseContext, mc *mail.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/** Read the request body **/
		type reqBody struct {
			AttendeeID string `json:"attendee_id"`
			PrivateKey string `json:"pubkey"` // hex-encoded private key
		}
		var body reqBody
		jsonDecoder := json.NewDecoder(r.Body)
		err := jsonDecoder.Decode(&body)
		if err != nil {
			errMsg := fmt.Sprintf("bad json request body: %s", err)
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}

		/** Attendee information + their their new private key **/
		privKey, err := util.Secp256k1FromHex(body.PrivateKey)
		if err != nil {
			errMsg := fmt.Sprintf("bad private key: %s", err)
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}
		attendeeAddress := util.IDToAddress(body.AttendeeID)

		/** Construct the secret for this and send the key transaction **/
		secret, commitment := util.CreateCommitment()
		err = mk.SendKeyTransaction(attendeeAddress, privKey.PubKey(), commitment)
		if err != nil {
			http.Error(w, "internal error. try again", http.StatusInternalServerError)
			return
		}

		/** Get the Attendee's email **/
		email, err := eb.EmailFromAttendeeID(body.AttendeeID)
		if err != nil {
			http.Error(w, "internal error. try again", http.StatusInternalServerError)
			return
		} else if len(email) == 0 {
			http.Error(w, "attendee id not present in the event", http.StatusNotFound)
			return
		}

		/** Store the private key **/
		ok := db.StoreKey(email, body.PrivateKey)
		if !ok {
			log.Error("could not store private key in the db")
		}

		/** Send the redirect **/
		err = mc.SendRedirectEmail(email, secret)
		if err != nil {
			http.Error(w, "email error. try again", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func keyGetter(db *models.DatabaseContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := mux.Vars(r)["email"]
		if !ok {
			http.Error(w, "email parameter required", http.StatusBadRequest)
		}

		key := db.GetKey(email)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(key))
	}
}
