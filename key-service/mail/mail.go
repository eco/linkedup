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
type Client interface {
	SendOnboardingEmail(string, sdk.AccAddress, string) error
	SendRecoveryEmail(string, sdk.AccAddress, string) error
}

type sesClient struct {
	ses *ses.SES
}


type mockClient struct {
}

// NewMockClient creates a mock email client session wrapper that just logs
// the template parameters so that the application can run locally without
// actually sending email
func NewMockClient() (client Client, err error) {
	client = mockClient{}
	return
}

// NewClient creates a new email client session wrapper
func NewClient(cfg client.ConfigProvider) (client Client, err error) {
	client = sesClient{
		ses: ses.New(cfg),
	}
	return
}

// SendOnboardingEmail will construct and send the email containing the initial
// onboarding message and URL with the given secret
func (c sesClient) SendOnboardingEmail(dest string, attendeeAddr sdk.AccAddress, secret string) error {
	redirectURI := makeRedirectURI(attendeeAddr, secret)

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
func (c sesClient) SendRecoveryEmail(dest string, attendeeAddr sdk.AccAddress, secret string) error {
	redirectURI := makeRedirectURI(attendeeAddr, secret)

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

func (c sesClient) sendEmailWithURL(dest string, url string, template string) (err error) {
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

func (c mockClient) SendOnboardingEmail(dest string, attendeeAddr sdk.AccAddress, secret string) error {
	redirectURI := makeRedirectURI(attendeeAddr, secret)

	log.Warnf("mock onboarding email with url: %s", redirectURI)
	return nil
}

func (c mockClient) SendRecoveryEmail(dest string, attendeeAddr sdk.AccAddress, secret string) error {
	redirectURI := makeRedirectURI(attendeeAddr, secret)

	log.Warnf("mock recovery email with url: %s", redirectURI)
	return nil
}

func makeRedirectURI(attendeeAddr sdk.AccAddress, secret string) string {
	return fmt.Sprintf(
		"https://linkedup.sfblockchainweek.io/claim?attendee=%s&secret=%s",
		attendeeAddr,
		secret,
	)
}
