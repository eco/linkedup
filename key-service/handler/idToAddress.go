package handler

import (
	"github.com/eco/longy/util"
	"github.com/gorilla/mux"
	"net/http"
)

// registerIDToAddress provides an deterministic converter from id to a
// cosmos address
func registerIDToAddress(r *mux.Router) {
	r.HandleFunc("/id/{id}", idToAddress()).Methods("GET")
}

func idToAddress() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		addr := util.IDToAddress(id).String()

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(addr))
	}
}
