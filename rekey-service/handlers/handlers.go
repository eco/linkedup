package handler

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

var log = logrus.WithField("module", "router")

func Router() http.Handler {
	r := mux.NewRouter()
	registerPing(r)

	return logRequests(r)
}

func logRequests(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("method=%s remote=%s url=%s content-length=%d",
			r.Method, r.RemoteAddr, r.URL, r.ContentLength)

		h.ServeHTTP(w, r)
	})
}
