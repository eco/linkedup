package mail

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	eb "github.com/eco/longy/key-service/eventbrite"
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
	SendOnboardingEmail(*eb.AttendeeProfile, sdk.AccAddress, string) error
	SendRecoveryEmail(*eb.AttendeeProfile, int, string) error
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
func (c sesClient) SendOnboardingEmail(profile *eb.AttendeeProfile, attendeeAddr sdk.AccAddress, secret string) error {
	redirectURI, err := makeOnboardingURI(profile, attendeeAddr, secret)

	if err != nil {
		log.Errorf(
			"unable to generate email URI: %s",
			err.Error(),
		)
		return err
	}

	log.Tracef("sending onboarding email to: %s", profile.Email)

	err = c.sendEmailWithURL(profile.Email, redirectURI, "linkedup-onboarding")

	if err != nil {
		log.Errorf(
			"unable to send onboarding email to %s: %s",
			profile.Email,
			err.Error(),
		)
	}

	return err
}

// SendRecoveryEmail will construct and send the email containing the account
// recovery message and URL with the given secret
func (c sesClient) SendRecoveryEmail(profile *eb.AttendeeProfile, id int, token string) error {
	redirectURI := makeRecoveryURI(id, token)

	log.Tracef("sending recovery email to: %s", profile.Email)

	err := c.sendEmailWithURL(profile.Email, redirectURI, "linkedup-rekey")

	if err != nil {
		log.Errorf(
			"unable to send recovery email to %s: %s",
			profile.Email,
			err.Error(),
		)
	}

	return nil
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

func (c mockClient) SendOnboardingEmail(profile *eb.AttendeeProfile, attendeeAddr sdk.AccAddress, secret string) error {
	redirectURI, err := makeOnboardingURI(profile, attendeeAddr, secret)

	if err != nil {
		return err
	}

	log.Warnf("mock onboarding email with url: %s", redirectURI)
	return nil
}

func (c mockClient) SendRecoveryEmail(profile *eb.AttendeeProfile, id int, token string) error {
	redirectURI := makeRecoveryURI(id, token)

	log.Warnf("mock recovery email with url: %s", redirectURI)
	return nil
}

func makeRecoveryURI(id int, token string) string {
	return fmt.Sprintf("http://localhost:5000/recover?id=%d&token=%s", id, token)
}

func makeOnboardingURI(profile *eb.AttendeeProfile, attendeeAddr sdk.AccAddress, secret string) (string, error) {
	jsonProfileData, err := json.Marshal(profile)
	if err != nil {
		log.WithError(err).Error("attendee profile serialization")
		return "", err
	}
	encodedProfileData := base64.URLEncoding.EncodeToString(jsonProfileData)

	return fmt.Sprintf(
		"http://localhost:5000/claim?attendee=%s&profile=%s&secret=%s",
		attendeeAddr,
		encodedProfileData,
		secret,
	), nil
}
