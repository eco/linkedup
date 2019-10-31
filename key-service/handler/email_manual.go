package handler

import (
	"encoding/json"
	"fmt"
	"github.com/badoux/checkmail"
	ebSession "github.com/eco/longy/key-service/eventbrite"
	"github.com/eco/longy/key-service/models"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
)

const (
	//EmailAuthHeader is the header name for getting the auth from a request
	EmailAuthHeader = "Authorization"
	//EmailAuthEnvKey the key name for auth token
	EmailAuthEnvKey = "EMAIL_AUTH"

	idKey = "id"
)

var authToken string

func init() {
	var isAuthSet bool
	authToken, isAuthSet = os.LookupEnv(EmailAuthEnvKey)
	if !isAuthSet {
		authToken = ""
		log.Errorf("environment variable %s not set!", EmailAuthEnvKey)
	} else {
		log.Info("email auth token set to env variable")
	}
}

func registerEmailManual(r *mux.Router, db *models.DatabaseContext, eb *ebSession.Session) {
	s := r.PathPrefix("/emails").Subrouter()
	s.HandleFunc("", setEmailForAttendee(db)).Methods(http.MethodPost, http.MethodOptions)
	s.HandleFunc(fmt.Sprintf("/{%s}", idKey), getEmailForAttendee(db, eb)).Methods(http.MethodGet, http.MethodOptions)
	s.Use(EmailAuthMiddleware)
}

func getEmailForAttendee(db *models.DatabaseContext, eb *ebSession.Session) func(
	http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars[idKey])
		if err != nil || id < 0 {
			http.Error(w, "id expected to be a positive integer", http.StatusBadRequest)
			return
		}

		attendee, found := eb.AttendeeProfile(id)
		if !found {
			http.Error(w, "unregistered badge", http.StatusNotFound)
			return
		}

		email := attendee.Email
		overrideEmail := db.GetEmail(id)

		type respBody struct {
			EventbriteEmail string `json:"eventbrite_email"`
			OverrideEmail   string `json:"override_email"`
		}

		resp := respBody{
			EventbriteEmail: email,
			OverrideEmail:   overrideEmail,
		}

		bz, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(bz)
	}
}

func setEmailForAttendee(db *models.DatabaseContext) func(
	http.ResponseWriter, *http.Request) {

	type emailBody struct {
		ID      int    `json:"id"`
		Address string `json:"address"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var eb emailBody
		if err := json.NewDecoder(r.Body).Decode(&eb); err != nil {
			http.Error(w, fmt.Sprintf("request body invalid: %s", err), http.StatusBadRequest)
			return
		} else if eb.ID < 0 {
			http.Error(w, "attendee id must be a positive integer", http.StatusBadRequest)
			return
		}

		if err := checkmail.ValidateFormat(eb.Address); err != nil {
			http.Error(w, "email address is empty or invalid", http.StatusBadRequest)
			return
		}

		ok := db.StoreEmail(eb.ID, eb.Address)

		if ok {
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

}

//EmailAuthMiddleware checks all requests
func EmailAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get(EmailAuthHeader)
		if auth != authToken {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
