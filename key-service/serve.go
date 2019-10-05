package rekeyservice

import (
	"context"
	"fmt"
	"github.com/eco/longy/key-service/eventbrite"
	"github.com/eco/longy/key-service/handler"
	"github.com/eco/longy/key-service/mail"
	"github.com/eco/longy/key-service/masterkey"
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
	ebSession  eventbrite.Session
	masterKey  masterkey.Key
	mailClient mail.Client
}

// NewService is the creator the the rekey-service
func NewService(ebSession eventbrite.Session, key masterkey.Key, mc mail.Client) Service {
	return Service{
		ebSession:  ebSession,
		masterKey:  key,
		mailClient: mc,
	}
}

// StartHTTP will block and start the http service binded on `port`
func (srv *Service) StartHTTP(port int) {
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler.Router(srv.ebSession, srv.masterKey, srv.mailClient),
	}

	// will block
	startServer(s)

	// server closed
	srv.Close()
}

// Close will release the resources used by the server
func (srv *Service) Close() {
	srv.mailClient.Close() //nolint
	log.Info("done")
}

func startServer(s *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Infof("listening on %s", s.Addr)
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-stop

	// graceful shutdown
	log.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.Shutdown(ctx) //nolint
}
