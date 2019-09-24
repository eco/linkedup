package rekeyservice

import (
	"context"
	"fmt"
	"github.com/eco/longy/rekey-service/eventbrite"
	"github.com/eco/longy/rekey-service/handlers"
	"github.com/eco/longy/rekey-service/masterkey"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var log = logrus.WithField("module", "rekeyservice")

// Service composes the required modules needed to manage the lifecycle
type Service struct {
	ebSession eventbrite.Session
	masterKey masterkey.Key
}

func NewService(ebSession eventbrite.Session, key masterkey.Key) Service {
	return Service{
		ebSession: ebSession,
		masterKey: key,
	}
}

// StartHTTP will block and start the http service binded on `port`
func (srv Service) StartHTTP(port int) {
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler.Router(),
	}

	startServer(s)
}

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