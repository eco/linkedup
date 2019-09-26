package mail

import (
	"crypto/tls"
	"encoding/hex"
	"fmt"
	gomail "github.com/go-mail/mail"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("module", "mail")

// Client used to send emails
type Client struct {
	sender gomail.SendCloser
}

// NewClient constructs `Client` with the corresponding credentials
func NewClient(host string, port int, username string, pwd string) (Client, error) {
	log.Infof("establishing connection with smtp server. %s:%d", host, port)
	d := gomail.NewDialer(host, port, username, pwd)
	d.TLSConfig = &tls.Config{
		ServerName: host,
	}
	d.StartTLSPolicy = gomail.MandatoryStartTLS

	// try and dial
	sender, err := d.Dial()
	if err != nil {
		return Client{}, fmt.Errorf("smtp conn: %s", err)
	}

	return Client{
		sender: sender,
	}, nil
}

// SendRekeyEmail will construct and send the email containing the redirect
// uri with the given signature
func (c Client) SendRekeyEmail(dest string, signature []byte) error {
	hexStr := hex.EncodeToString(signature)

	// construct message
	m := gomail.NewMessage()
	m.SetHeader("From", "testecolongy@gmail.com")
	m.SetHeader("To", dest)
	m.SetHeader("From", "alex@example.com")
	m.SetHeader("Subject", "Reset keys and re-enter the longy game")
	m.SetBody("text/html", fmt.Sprintf("<b>Hello!</b> Signature: %s", hexStr))

	err := gomail.Send(c.sender, m)
	if err != nil {
		log.WithError(err).Warnf("failed email delivery. email: %s", dest)
	}
	return err
}

// Close will terminate the connnection with the smtp server
func (c Client) Close() error {
	log.Info("terminating connection with smtp server")
	return c.sender.Close()
}
