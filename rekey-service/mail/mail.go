package mail

import (
	"crypto/tls"
	"fmt"
	gomail "github.com/go-mail/mail"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("module", "mail")

// Client used to sent emails
type Client struct {
	sender gomail.SendCloser
}

func NewClient(host string, port int, username string, pwd string) (Client, error) {
	log.Info("establishing connection with smtp server")
	d := gomail.NewDialer(host, port, username, pwd)
	// TODO: change this for production!
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
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
	// construct message
	m := gomail.NewMessage()
	m.SetHeader("To", dest)
	m.SetHeader("From", "alex@example.com")
	m.SetHeader("Subject", "Reset keys and re-enter the longy game")
	m.SetBody("text/html", "<b>Hello!</b>")

	return gomail.Send(c.sender, m)
}

func (c Client) Close() error {
	return c.sender.Close()
}
