package handler

import (
	"github.com/eco/longy/key-service/types"
	"github.com/eco/longy/key-service/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

// Router returns the root http Handler
func Router(srv types.Service) http.Handler {
	r := mux.NewRouter()
	registerPing(r)
	registerRekey(r, srv)

	return middleware.LogHTTP(r)
}
