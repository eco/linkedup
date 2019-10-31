package handler

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/eventbrite"
	ebSession "github.com/eco/longy/key-service/eventbrite"
	longyClnt "github.com/eco/longy/key-service/longyclient"
	"github.com/eco/longy/key-service/mail"
	"github.com/eco/longy/key-service/masterkey"
	"github.com/eco/longy/key-service/models"
	"github.com/eco/longy/util"
	"github.com/gorilla/mux"
)

// AttendeeInfo -
type AttendeeInfo struct {
	Address sdk.AccAddress              `json:"address"`
	Profile *eventbrite.AttendeeProfile `json:"attendee"`

	// private key information
	CosmosPrivateKey string `json:"cosmos_private_key"`
	RSAPrivateKey    string `json:"rsa_private_key"`
	RSAPublicKey     string `json:"rsa_public_key"`

	// needed to claim the attendee account
	CommitmentSecret string          `json:"commitment_secret"`
	Commitment       util.Commitment `json:"commitment"`

	ImageUploadURL string `json:"image_upload_url"`
}

func registerKey(
	r *mux.Router,
	eb *ebSession.Session,
	mk *masterkey.MasterKey,
	db *models.DatabaseContext,
	mc mail.Client) {

	// POST
	r.HandleFunc("/key", key(eb, mk, db, mc)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/recover", keyRecover(db, mk, mc)).Methods(http.MethodPost, http.MethodOptions)

	// GET
	r.HandleFunc("/recover/{id}/{token}", keyRetrieval(db)).Methods(http.MethodGet, http.MethodOptions)
}

// All core logic is implemented here. If there are plans to expand this service,
// logic (email retrieval, etc) can be lifted into http middleware to allow for better
// composability
//nolint: gocyclo, gocritic
func key(eb *ebSession.Session,
	mk *masterkey.MasterKey,
	db *models.DatabaseContext,
	mc mail.Client) http.HandlerFunc {
	type reqBody struct {
		AttendeeID int `json:"attendee_id"`

		// private key information
		CosmosPrivateKey string `json:"cosmos_private_key"`
		RSAPrivateKey    string `json:"rsa_private_key"`
		RSAPublicKey     string `json:"rsa_public_key"`

		UseVerification bool `json:"use_verification"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		/** Read the request body **/
		var body reqBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, fmt.Sprintf("request body: %s", err), http.StatusBadRequest)
			return
		} else if body.AttendeeID < 0 {
			http.Error(w, "attendee id must be a positive integer", http.StatusBadRequest)
			return
		} else if body.CosmosPrivateKey == "" || body.RSAPrivateKey == "" || body.RSAPublicKey == "" {
			http.Error(w, "cosmos private key & rsa private/public keys must be present", http.StatusBadRequest)
			return
		}

		/** Check if this attendee is in eventbrite **/
		profile, found := eb.AttendeeProfile(body.AttendeeID)
		if !found {
			http.Error(w, "non-registered badge id", http.StatusNotFound)
			return
		}

		/** Check if this attendee already has info registered **/
		infoBz, err := db.GetAttendeeInfo(body.AttendeeID)
		if err != nil {
			http.Error(w, "key-service down", http.StatusServiceUnavailable)
			return
		} else if len(infoBz) != 0 {
			http.Error(w, "attendee info onboarded. use /recover instead", http.StatusConflict)
			return
		}

		/** Sanity check on the private key **/
		_, err = util.Secp256k1FromHex(body.CosmosPrivateKey)
		if err != nil {
			http.Error(w, fmt.Sprintf("bad cosmos private key: %s", err), http.StatusBadRequest)
			return
		}

		/** Store the attendee information **/

		imageUploadURL, err := db.GetImageUploadURL(body.AttendeeID)
		if err != nil {
			http.Error(w, "failed to sign image upload URL", http.StatusServiceUnavailable)
			return
		}

		// generate the unique <secret / commitment> pair for this attendee
		secret, commitment := util.CreateCommitment()
		info := &AttendeeInfo{
			Address: util.IDToAddress(fmt.Sprintf("%d", body.AttendeeID)),
			Profile: profile,

			CosmosPrivateKey: body.CosmosPrivateKey,
			RSAPrivateKey:    body.RSAPrivateKey,
			RSAPublicKey:     body.RSAPublicKey,

			CommitmentSecret: secret,
			Commitment:       commitment,

			ImageUploadURL: imageUploadURL,
		}
		bz, err := json.Marshal(info)
		if err != nil {
			log.WithError(err).WithField("data", info).Error("marshaling attendee info")
			http.Error(w, "key storage service down", http.StatusInternalServerError)
			return
		}
		if ok := db.StoreAttendeeInfo(body.AttendeeID, bz); !ok {
			http.Error(w, "key storage service down", http.StatusServiceUnavailable)
			return
		}

		keyed := false // attendee information stored for the first time. Account could not have been keyed prior
		keyAndEmail(mk, db, mc,
			body.AttendeeID, info, true, keyed, body.UseVerification)(w, r)
	}
}

func keyRecover(
	db *models.DatabaseContext,
	mk *masterkey.MasterKey,
	mc mail.Client) http.HandlerFunc {
	type reqBody struct {
		AttendeeID      int  `json:"attendee_id"`
		UseVerification bool `json:"use_verification"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		/** Read request body**/
		var body reqBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, fmt.Sprintf("request body: %s", err), http.StatusBadRequest)
			return
		} else if body.AttendeeID < 0 {
			http.Error(w, "attendee id must be a positive integer", http.StatusBadRequest)
			return
		}

		/** Retrieve attendee information from the database **/
		infoBz, err := db.GetAttendeeInfo(body.AttendeeID)
		if err != nil {
			http.Error(w, "key-service down", http.StatusServiceUnavailable)
			return
		} else if len(infoBz) == 0 {
			http.Error(w, "attendee not found", http.StatusNotFound)
			return
		}
		var attendeeInfo AttendeeInfo
		err = json.Unmarshal(infoBz, &attendeeInfo)
		if err != nil {
			http.Error(w, "corrupt attendee information", http.StatusInternalServerError)
			return
		}

		keyed, err := longyClnt.IsAttendeeKeyed(body.AttendeeID)
		if err != nil {
			http.Error(w, "key-service down", http.StatusServiceUnavailable)
			return
		}

		keyAndEmail(mk, db, mc,
			body.AttendeeID, &attendeeInfo, false, keyed, body.UseVerification)(w, r)
	}
}

// helper used by both `key` and `keyRecover`
//nolint:gocyclo
func keyAndEmail(
	mk *masterkey.MasterKey,
	db *models.DatabaseContext,
	mc mail.Client,

	id int,
	info *AttendeeInfo,
	onboarding, keyed, useVerification bool,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := fmt.Sprintf("%d", id)

		// key the attendee if required
		if !keyed {
			privKey, err := util.Secp256k1FromHex(info.CosmosPrivateKey)
			if err != nil {
				http.Error(w, "corrupt attendee private key", http.StatusInternalServerError)
				return
			}
			err = mk.SendKeyTransaction(util.IDToAddress(idStr), privKey.PubKey(), info.Commitment)
			if err != nil && err != masterkey.ErrAlreadyKeyed {
				http.Error(w, "internal error. try again", http.StatusInternalServerError)
				return
			}
		}

		// unique token to retrieve stored info. Store this token if we are `useVerification` or `!onboarding`.
		// !onboarding indicates this was instantiated via recovery
		token := generateVerificationToken()
		if useVerification || !onboarding {
			if ok := db.StoreVerificationToken(id, token); !ok {
				http.Error(w, "key-service down", http.StatusServiceUnavailable)
				return
			}
		}

		// send the email
		var err error
		//check if the email has been changed manually for the attendee
		storedEmail := db.GetEmail(info.Profile.ID)
		if storedEmail != "" {
			info.Profile.Email = storedEmail
		}
		if useVerification {
			err = mc.SendVerificationEmail(db, info.Profile.Email, token)
		} else {
			if onboarding {
				// onboarding email
				err = mc.SendOnboardingEmail(db, info.Address, info.Profile, info.CommitmentSecret, info.ImageUploadURL)
			} else {
				// recovery email
				err = mc.SendRecoveryEmail(db, info.Profile, idStr, token)
			}
		}

		if err != nil {
			http.Error(w, "email error. try again", http.StatusServiceUnavailable)
			return
		}

		maskedEmail := maskEmail(info.Profile.Email)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(maskedEmail)) //nolint
	}
}

// retrieve attendee information with the given verification token
func keyRetrieval(db *models.DatabaseContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		token := vars["token"]
		id, err := strconv.Atoi(vars["id"])
		if err != nil || id < 0 {
			http.Error(w, "id expected to be a positive integer", http.StatusBadRequest)
			return
		}

		expectedToken, err := db.GetVerificationToken(id)
		if err != nil {
			http.Error(w, "key-service down", http.StatusServiceUnavailable)
			return
		}
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

/** Helpers **/

// used to generate a verification code 6 digits in length
var table = [10]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func generateVerificationToken() string {
	b := make([]byte, 6)
	n, err := io.ReadAtLeast(rand.Reader, b, 6)
	if n != 6 {
		panic(err)
	}

	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func maskEmail(email string) string {
	splitEmail := strings.Split(email, "@")
	if len(splitEmail) != 2 {
		return "malformed email. see support"
	}

	username := splitEmail[0]

	var maskedUsername string
	switch {
	case len(username) == 1:
		// we won't mask this
		maskedUsername = username
	case len(username) == 2:
		maskedUsername = string(username[0]) + "*"
	default:
		maskedUsername = string(username[0]) + strings.Repeat("*", len(username)-2) + string(username[len(username)-1])
	}

	maskedEmail := maskedUsername + "@" + splitEmail[1]
	return maskedEmail
}
