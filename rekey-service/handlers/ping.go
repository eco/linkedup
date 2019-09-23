package handler

import (
	"github.com/gorilla/mux"
	"net/http"
)

var _ http.HandlerFunc = ping

func registerPing(r *mux.Router) {
	r.HandleFunc("/ping", ping).Methods("GET")
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
	return
}
