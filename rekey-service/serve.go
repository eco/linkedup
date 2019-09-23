package rekeyservice

import (
	"context"
	"fmt"
	"github.com/eco/longy/rekey-service/handlers"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var log = logrus.WithField("module", "rekeyservice")

// StartHttpService will block and start the http service binded on `port`
func StartHttpService(port int) {
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler.Router(),
	}

	startServer(s)
}

// TODO: https server?

func startServer(s *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Infof("listening on %s", s.Addr)
		if err := s.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-stop

	// graceful shutdown
	log.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.Shutdown(ctx)
}
