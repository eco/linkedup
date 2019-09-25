package mail

import (
	"github.com/sirupsen/logrus"
	//"net/smtp"
)

var log = logrus.WithField("module", "mail")

// Client used to sent emails
type Client struct {
}

func NewClient() (Client, error) {
	log.Info("establishing connection with smtp server")
	return Client{}, nil
}

// SendRekeyEmail will construct and send the email containing the redirect
// uri with the given signature
func (c Client) SendRekeyEmail(dest string, signature []byte) error {
	return nil
}
