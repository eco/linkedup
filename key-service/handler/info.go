package handler

import (
	"encoding/json"
	"fmt"
	"github.com/eco/longy/key-service/mail"
	"github.com/eco/longy/key-service/models"
	"github.com/gorilla/mux"
	"net/http"
)

func registerInfo(r *mux.Router, db *models.DatabaseContext, mc mail.Client) {
	r.HandleFunc("/sendEmail", sendEmailToAttendee(db, mc)).Methods(http.MethodPost, http.MethodOptions)
}

//nolint:gocyclo
func sendEmailToAttendee(db *models.DatabaseContext, mc mail.Client) func(
	http.ResponseWriter, *http.Request) {
	type sendBody struct {
		ID    int    `json:"id"`
		Token string `json:"token"`
		Data  string `json:"data"`
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

		expectedToken, err := db.GetVerificationToken(info.Profile.ID)
		if err != nil {
			http.Error(w, "could not retrieve expected token for account", http.StatusServiceUnavailable)
			return
		}

		if len(expectedToken) == 0 {
			http.Error(w, "attendee has no auth token stored", http.StatusUnauthorized)
			return
		}

		if expectedToken != sb.Token {
			http.Error(w, "incorrect auth token", http.StatusUnauthorized)
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
