package rekeyservice

import (
	"context"
	"fmt"
	"github.com/eco/longy/rekey-service/eventbrite"
	"github.com/eco/longy/rekey-service/handler"
	"github.com/eco/longy/rekey-service/mail"
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
	ebSession  eventbrite.Session
	masterKey  masterkey.Key
	mailClient mail.Client
}

func NewService(ebSession eventbrite.Session, key masterkey.Key, mc mail.Client) Service {
	return Service{
		ebSession:  ebSession,
		masterKey:  key,
		mailClient: mc,
	}
}

// StartHTTP will block and start the http service binded on `port`
func (srv Service) StartHTTP(port int) {
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler.Router(srv.ebSession, srv.masterKey, srv.mailClient),
	}

	// will block
	startServer(s)

	srv.Close()
}

func (srv Service) Close() error {
	err := srv.mailClient.Close()
	log.Info("done")
	return err
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

	s.Shutdown(ctx)
}
