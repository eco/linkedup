package handler

import (
	"encoding/hex"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type EmailRetrievalService interface {
	EmailFromAttendeeID(id int) (string, error)
}

type MasterKey interface {
	RekeySignature(id, nonce int) ([]byte, error)
}

type MailService interface {
	SendRekeyEmail(email string, sig []byte) error
}

func registerRekey(r *mux.Router, eb EmailRetrievalService, mk MasterKey, mc MailService) {
	r.HandleFunc("/rekey/{id:[0-9]+}", rekey(eb, mk, mc)).Methods("GET")
}

// All core logic is implemented here. If there are plans to expand this service,
// logic (email retrieval, etc) can be lifted into http middleware to allow for better
// composability
func rekey(eb EmailRetrievalService, mk MasterKey, mc MailService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nonceStr := r.FormValue("nonce")
		if len(nonceStr) == 0 {
			http.Error(w, "nonce expected as a url parameter", http.StatusBadRequest)
			return
		}

		// retrieve id & nonce from the url path
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "attendee id must be in decimal format", http.StatusBadRequest)
			return
		}
		nonce, err := strconv.Atoi(nonceStr)
		if err != nil {
			http.Error(w, "nonce must be in decimal format", http.StatusBadRequest)
			return
		}

		sig, err := mk.RekeySignature(id, nonce)
		if err != nil {
			http.Error(w, "internal error. try again", http.StatusInternalServerError)
			return
		}
		hexStr := hex.EncodeToString(sig)

		email, err := eb.EmailFromAttendeeID(id)
		if err != nil {
			http.Error(w, "internal error. try again", http.StatusInternalServerError)
			return
		}

		err = mc.SendRekeyEmail(email, sig)
		if err != nil {
			var body map[string]string
			body["msg"] = "unable to send email"
			body["signature"] = hexStr
			data, err := json.Marshal(body)
			if err != nil {
				log.WithField("handler", "rekey").Error("json parsing")
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusAccepted)
			w.Write(data)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(hexStr))
	}
}
