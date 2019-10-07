package mail

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/ses"
)

var (
	log = logrus.WithField("module", "mail")

	gmEmail = "LinkedUp Game <gm@linkedup.sfblockchainweek.io>"
)

// Client used to send emails
type Client struct {
	ses *ses.SES
}

// NewClient creates a new email client session wrapper
func NewClient(cfg client.ConfigProvider) (client Client, err error) {
	client = Client{
		ses: ses.New(cfg),
	}
	return
}

// SendOnboardingEmail will construct and send the email containing the initial
// onboarding message and URL with the given secret
func (c Client) SendOnboardingEmail(dest string, attendeeAddr sdk.AccAddress, secret string) error {
	redirectURI := fmt.Sprintf("https://linkedup.sfblockchainweek.io/claim?attendee=%s&secret=%s", attendeeAddr, secret)

	log.Tracef("sending onboarding email to: %s", dest)

	err := c.sendEmailWithURL(dest, redirectURI, "linkedup-onboarding")

	if err != nil {
		log.Errorf(
			"unable to send onboarding email to %s: %s",
			dest,
			err.Error(),
		)
	}

	return err
}

// SendRecoveryEmail will construct and send the email containing the account
// recovery message and URL with the given secret
func (c Client) SendRecoveryEmail(dest string, attendeeAddr sdk.AccAddress, secret string) error {
	redirectURI := fmt.Sprintf("https://linkedup.sfblockchainweek.io/claim?attendee=%s&secret=%s", attendeeAddr, secret)

	log.Tracef("sending recovery email to: %s", dest)

	err := c.sendEmailWithURL(dest, redirectURI, "linkedup-rekey")

	if err != nil {
		log.Errorf(
			"unable to send recovery email to %s: %s",
			dest,
			err.Error(),
		)
	}

	return err
}

func (c Client) sendEmailWithURL(dest string, url string, template string) (err error) {
	templateData := fmt.Sprintf("{\"url\":\"%s\"}", url)

	_, err = c.ses.SendTemplatedEmail(&ses.SendTemplatedEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{&dest},
		},
		Source: &gmEmail,
		Template: &template,
		TemplateData: &templateData,
	})
	return
}
