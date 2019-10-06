package types

import (
	"github.com/eco/longy/key-service/eventbrite"
	"github.com/eco/longy/key-service/mail"
	"github.com/eco/longy/key-service/masterkey"
)

// Service defines the state of a server that provides a service
type Service interface {
	// MailClient provides access to the mail client resource handle
	MailClient() *mail.Client

	// Eventbrite provides access to the Eventbrite API session
	Eventbrite() *eventbrite.Session

	// MasterKey provides access to the signing key that allows access to the
	// chain's administrative functions
	MasterKey() *masterkey.Key
}

