package handler

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/key-service/mail"
	"github.com/eco/longy/key-service/models"
	"github.com/eco/longy/x/longy/crypto"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"net/http"
)

func registerInfo(r *mux.Router, db *models.DatabaseContext, mc mail.Client) {
	r.HandleFunc("/sendEmail", sendEmailToAttendee(db, mc)).Methods(http.MethodPost, http.MethodOptions)
}

func sendEmailToAttendee(db *models.DatabaseContext, mc mail.Client) func(
	http.ResponseWriter, *http.Request) {
	type sendBody struct {
		ID   int    `json:"id"`
		Sig  string `json:"sig"`
		Data string `json:"data"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var sb sendBody
		if err := json.NewDecoder(r.Body).Decode(&sb); err != nil {
			http.Error(w, fmt.Sprintf("request body invalid: %s", err), http.StatusBadRequest)
			return
		} else if sb.ID < 0 {
			http.Error(w, "attendee id must be a positive integer", http.StatusBadRequest)
			return
		}

		infoBz, err := db.GetAttendeeInfo(sb.ID)
		if err != nil {
			http.Error(w, "key-service down", http.StatusServiceUnavailable)
			return
		} else if len(infoBz) == 0 {
			http.Error(w, "attendee does not exist", http.StatusNotFound)
			return
		}

		var info AttendeeInfo
		err = json.Unmarshal(infoBz, &info)
		if err != nil {
			http.Error(w, "failed to unmarshal attendee info", http.StatusInternalServerError)
			return
		}

		err = validSig(info.CosmosPrivateKey, sb.Sig)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		//check if the email has been changed manually for the attendee
		storedEmail := db.GetEmail(info.Profile.ID)
		if storedEmail != "" {
			info.Profile.Email = storedEmail
		}

		err = mc.SendAttendeeSharedInfoEmail(db, info.Profile.Email, sb.Data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		maskedEmail := maskEmail(info.Profile.Email)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(maskedEmail))

	}
}

func validSig(privKey string, sig string) error {
	var priv secp256k1.PrivKeySecp256k1
	tmp := []byte(privKey)
	copy(priv[:], tmp)
	addrString := sdk.AccAddress(priv.PubKey().Address()).String()
	return crypto.ValidateSig(priv.PubKey(), addrString, sig)
}
