package middleware

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var log = logrus.WithField("module", "middleware")

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

// LogHTTP will create and info level entry about the incoming request
func LogHTTP(underlying http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sw := statusWriter{
			ResponseWriter: w,
			status:         0,
			length:         0,
		}

		start := time.Now()
		underlying.ServeHTTP(&sw, r)
		latency := time.Since(start)

		log.Infof("method=%s, remote=%s, url=%s, content-length=%d, status=%d, latency=%s",
			r.Method, r.RemoteAddr, r.URL, r.ContentLength, sw.status, latency)
	}
}
