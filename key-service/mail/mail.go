package mail

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	eb "github.com/eco/longy/key-service/eventbrite"
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

// SendRedirectEmail will construct and send the email containing the redirect
// uri with the given secret and attendee profile
func (c Client) SendRedirectEmail(profile *eb.AttendeeProfile, attendeeAddr sdk.AccAddress, secret string) error {
	jsonProfileData, err := json.Marshal(profile)
	if err != nil {
		log.WithError(err).Error("attendee profile serialization")
		return err
	}
	encodedProfileData := base64.StdEncoding.EncodeToString(jsonProfileData)

	redirectURI := fmt.Sprintf("http://longygame.com/claim?attendee=%s&profile=%s,secret=%s",
		attendeeAddr, encodedProfileData, secret)

	// construct message
	m := gomail.NewMessage()
	m.SetHeader("From", "testecolongy@gmail.com")
	m.SetHeader("To", profile.Email)
	m.SetHeader("From", "alex@example.com")
	m.SetHeader("Subject", "Onboard to the the longy game")
	m.SetBody("text/html", fmt.Sprintf("<b>Hello!</b> enter the game -> %s", redirectURI))

	if err := gomail.Send(c.sender, m); err != nil {
		log.WithError(err).WithField("dest", profile.Email).
			Error("failed email delivery")

		return err
	}

	return nil
}

// SendRecoveryEmail will send an email with the `authToken` required to hit the /recover/{id}/{token} endpoint and retrieve
// the keys that are stored in the backend
func (c Client) SendRecoveryEmail(dest, authToken string, id int) error {
	redirectURI := fmt.Sprintf("http://longygame.com/recover?id=%d&token=%s", id, authToken)

	// construct message
	m := gomail.NewMessage()
	m.SetHeader("From", "testecolongy@gmail.com")
	m.SetHeader("To", dest)
	m.SetHeader("From", "alex@example.com")
	m.SetHeader("Subject", "Onboard to the the longy game")
	m.SetBody("text/html", fmt.Sprintf("<b>Hello!</b> recover your account -> %s", redirectURI))

	if err := gomail.Send(c.sender, m); err != nil {
		log.WithError(err).WithField("dest", dest).
			Error("failed email delivery")

		return err
	}
	return nil
}

// Close will terminate the connnection with the smtp server
func (c Client) Close() error {
	log.Info("terminating connection with smtp server")
	return c.sender.Close()
}
