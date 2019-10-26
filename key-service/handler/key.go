package handler

import (
	"encoding/json"
	"fmt"
	"github.com/eco/longy/eventbrite"
	ebSession "github.com/eco/longy/key-service/eventbrite"
	"github.com/eco/longy/key-service/mail"
	"github.com/eco/longy/key-service/masterkey"
	"github.com/eco/longy/key-service/models"
	"github.com/eco/longy/util"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

// AttendeeInfo -
type AttendeeInfo struct {
	Profile *eventbrite.AttendeeProfile `json:"attendee"`

	// private key information
	CosmosPrivateKey string `json:"cosmos_private_key"`
	RSAPrivateKey    string `json:"RSA_key"`
}

func registerKey(
	r *mux.Router,
	eb *ebSession.Session,
	mk *masterkey.MasterKey,
	db *models.DatabaseContext,
	mc mail.Client) {

	r.HandleFunc("/key", key(eb, mk, db, mc)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/recover", keyRecover(db, eb, mc)).Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc("/recover/{id}/{token}", keyGetter(db)).Methods(http.MethodGet, http.MethodOptions)
}

// All core logic is implemented here. If there are plans to expand this service,
// logic (email retrieval, etc) can be lifted into http middleware to allow for better
// composability
//nolint: gocyclo
func key(eb *ebSession.Session,
	mk *masterkey.MasterKey,
	db *models.DatabaseContext,
	mc mail.Client) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		/** Read the request body **/
		type reqBody struct {
			AttendeeID int `json:"attendee_id"`

			// private key information
			CosmosPrivateKey string `json:"cosmos_private_key"`
			RSAPrivateKey    string `json:"rsa_private_key"`
		}
		var body reqBody
		jsonDecoder := json.NewDecoder(r.Body)
		err := jsonDecoder.Decode(&body)
		if err != nil {
			errMsg := fmt.Sprintf("bad json request body: %s", err)
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}

		/*
			Check if this attendee already has info registered
			hasInfo, err := db.HasAttendeeInfo(body.AttendeeID)
			if err != nil {
				http.Error(w, "key-service down", http.StatusServiceUnavailable)
				return
			} else if hasInfo {
				http.Error(w, "attendee info already stored", http.StatusConflict)
				return
			}
		*/

		/** Attendee information + their their new private key **/
		privKey, err := util.Secp256k1FromHex(body.CosmosPrivateKey)
		if err != nil {
			errMsg := fmt.Sprintf("bad cosmos private key: %s", err)
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}
		attendeeAddress := util.IDToAddress(fmt.Sprintf("%d", body.AttendeeID))

		/** Get the Attendee's profile **/
		profile, ok := eb.AttendeeProfile(body.AttendeeID)
		if !ok {
			http.Error(w, "key-service down", http.StatusServiceUnavailable)
			return
		}

		/** Store the attendee information **/
		info := AttendeeInfo{
			Profile:          profile,
			CosmosPrivateKey: body.CosmosPrivateKey,
			RSAPrivateKey:    body.RSAPrivateKey,
		}
		bz, err := json.Marshal(info)
		if err != nil {
			log.WithError(err).WithField("data", info).
				Error("marshaling attendee info")
			// this is a server side error that should be covered. 500
			http.Error(w, "key storage service down", http.StatusInternalServerError)
			return
		}
		ok = db.StoreAttendeeInfo(body.AttendeeID, bz)
		if !ok {
			http.Error(w, "key storage service down", http.StatusServiceUnavailable)
			return
		}

		/** Construct the secret for this and send the key transaction **/
		secret, commitment := util.CreateCommitment()
		err = mk.SendKeyTransaction(attendeeAddress, privKey.PubKey(), commitment)
		if err != nil {
			if err == masterkey.ErrAlreadyKeyed {
				http.Error(w, "id has already been keyed", http.StatusUnauthorized)
			} else {
				http.Error(w, "internal error. try again", http.StatusInternalServerError)
			}
			return
		}
		log.Infof("keyed badge id: %d", body.AttendeeID)

		imageUploadURL, err := db.GetImageUploadURL(strconv.Itoa(body.AttendeeID))
		if err != nil {
			http.Error(w, "failed to sign image upload URL", http.StatusInternalServerError)
			return
		}

		/** Send the redirect **/
		err = mc.SendOnboardingEmail(profile, attendeeAddress, secret, imageUploadURL)
		if err != nil {
			http.Error(w, "email error. try again", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func keyRecover(db *models.DatabaseContext, eb *ebSession.Session, mc mail.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/** Read the attendee id from the body **/
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "unable to read request body", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(string(body))
		if err != nil || id < 0 {
			http.Error(w, "body expected to be a positive integer denoting the attendee id", http.StatusBadRequest)
			return
		}

		/** Retrieve email information **/
		profile, ok := eb.AttendeeProfile(id)
		if !ok {
			http.Error(w, "attendee not found", http.StatusNotFound)
			return
		}

		/** Create an auth token and send recovery email **/
		token := uuid.New().String()
		if ok := db.StoreAuthToken(id, token); !ok {
			http.Error(w, "key service down", http.StatusServiceUnavailable)
			return
		}

		if err := mc.SendRecoveryEmail(profile, id, token); err != nil {
			http.Error(w, "key service down", http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func keyGetter(db *models.DatabaseContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		token := vars["token"]

		id, err := strconv.Atoi(vars["id"])
		if err != nil || id < 0 {
			http.Error(w, "id expected to be a positive integer", http.StatusBadRequest)
			return
		}

		expectedToken := db.GetAuthToken(id)
		if len(expectedToken) == 0 {
			http.Error(w, "attendee has not attempted recovery", http.StatusUnauthorized)
			return
		} else if expectedToken != token {
			http.Error(w, "incorrect auth token", http.StatusUnauthorized)
			return
		}

		// authentication checks passed
		bz, err := db.GetAttendeeInfo(id)
		switch {
		case err != nil:
			http.Error(w, "key-service down", http.StatusServiceUnavailable)
		case bz == nil:
			http.Error(w, "not found", http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusOK)
			w.Write(bz) //nolint
		}
	}
}
