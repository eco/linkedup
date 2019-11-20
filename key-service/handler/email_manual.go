package handler

import (
	"encoding/json"
	"fmt"
	"github.com/badoux/checkmail"
	ebSession "github.com/eco/longy/key-service/eventbrite"
	"github.com/eco/longy/key-service/mail"
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

func registerEmailManual(r *mux.Router, db *models.DatabaseContext, eb *ebSession.Session, mc mail.Client) {
	s := r.PathPrefix("/emails").Subrouter()
	s.HandleFunc("", setEmailForAttendee(db)).Methods(http.MethodPost, http.MethodOptions)
	s.HandleFunc(fmt.Sprintf("/{%s}", idKey), getEmailForAttendee(db, eb)).Methods(http.MethodGet, http.MethodOptions)
	s.HandleFunc("/sendReceiveInfo", sendReceiveInfo(db, eb, mc)).Methods(http.MethodPost, http.MethodOptions)
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

//sendReceiveInfo send emails to all the attendees that have keyed their accounts and onboarded
func sendReceiveInfo(db *models.DatabaseContext, eb *ebSession.Session, mc mail.Client) func(
	http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		//get keyed accounts
		attendees := eb.GetAttendees()
		if len(attendees) == 0 {
			http.Error(w, fmt.Sprintf("there are no attendees"), http.StatusBadRequest)
		}

		//generate recover tokens for them
		var ids []int
		eco := make(map[int]interface{}, 2)
		eco[1454663009] = "stoyan"
		//eco[1284763465] = "andy"
		//for k := range attendees {
		for k := range eco {
			infoBz, err := db.GetAttendeeInfo(k)
			if err != nil || len(infoBz) == 0 {
				continue
			}

			var info AttendeeInfo
			err = json.Unmarshal(infoBz, &info)
			if err != nil {
				continue
			}

			token := generateVerificationToken()
			if ok := db.StoreVerificationToken(k, token); !ok {
				fmt.Println("error storing verification token")
				continue
			}

			//check if the email has been changed manually for the attendee
			storedEmail := db.GetEmail(info.Profile.ID)
			if storedEmail != "" {
				info.Profile.Email = storedEmail
			}

			err = mc.SendExportEmail(db, info.Profile.Email, k, token)
			if err != nil {
				fmt.Printf("failed to send export email to  : %s", info.Profile.Email)
				continue
			}
			ids = append(ids, info.Profile.ID)
		}
		type res struct {
			Ids []int `json:"ids"`
		}
		bz, err := json.Marshal(res{ids})
		if err != nil {
			w.WriteHeader(http.StatusOK)
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
