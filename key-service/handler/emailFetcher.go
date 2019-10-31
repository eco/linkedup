package handler

import (
	"net/http"
	"strconv"
	"strings"

	ebSession "github.com/eco/longy/key-service/eventbrite"
	"github.com/eco/longy/key-service/models"
	"github.com/gorilla/mux"
)

//nolint
func registerEmailFetcher(r *mux.Router, db *models.DatabaseContext, eb *ebSession.Session) {
	r.HandleFunc("/email/{id:[0-9]+}", getEmailForID(db, eb)).Methods(http.MethodPost, http.MethodOptions)
}

//nolint
func getEmailForID(db *models.DatabaseContext, eb *ebSession.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil || id < 0 {
			http.Error(w, "id must be a positive decimal number", http.StatusBadRequest)
			return
		}

		var email string
		email = db.GetEmail(id)
		if email != "" {
			profile, found := eb.AttendeeProfile(id)
			if !found {
				http.Error(w, "non-registered badge id", http.StatusNotFound)
				return
			} else if profile.Email == "" {
				http.Error(w, "missing email from attendee profile", http.StatusConflict)
				return
			}

			email = profile.Email
		}

		splitEmail := strings.Split(email, "@")
		if len(splitEmail) != 2 {
			http.Error(w, "malformed email", http.StatusConflict)
			return
		}

		username := splitEmail[0]

		var maskedUsername string
		if len(username) == 1 {
			maskedUsername = username
		} else if len(username) == 2 {
			maskedUsername = string(username[0]) + "*"
		} else {
			maskedUsername = string(username[0]) + strings.Repeat("*", len(username)-2) + string(username[len(username)-1])
		}

		maskedEmail := maskedUsername + "@" + splitEmail[1]
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(maskedEmail))
	}
}
