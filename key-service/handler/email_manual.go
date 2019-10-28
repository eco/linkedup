package handler

import (
	"fmt"
	"github.com/eco/longy/key-service/eventbrite"
	"github.com/eco/longy/key-service/models"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

const (
	//EmailAuthHeader is the header name for getting the auth from a request
	EmailAuthHeader = "Authorization"
	//EmailAuthEnvKey the key name for auth token
	EmailAuthEnvKey = "EMAIL_AUTH"
)

var authToken string

func init() {
	var isAuthSet bool
	authToken, isAuthSet = os.LookupEnv(EmailAuthEnvKey)
	if !isAuthSet {
		panic("auth token for email not set, set EMAIL_AUTH in env vars")
	} else {
		fmt.Println("email auth token set to env variable")
	}
}

func registerEmailManual(r *mux.Router, eb *eventbrite.Session, db *models.DatabaseContext) {
	s := r.PathPrefix("/emails").Subrouter()
	s.HandleFunc("", setEmailForAttendee(eb, db)).Methods(http.MethodPost, http.MethodOptions)
	s.HandleFunc("/{id}", getEmailForAttendee(eb, db)).Methods(http.MethodGet, http.MethodOptions)
	s.Use(EmailAuthMiddleware)
}

func getEmailForAttendee(session *eventbrite.Session, context *models.DatabaseContext) func(
	http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}

func setEmailForAttendee(session *eventbrite.Session, context *models.DatabaseContext) func(
	http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {

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
