package mail

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	eb "github.com/eco/longy/eventbrite"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/ses"
)

var (
	log = logrus.WithField("module", "mail")

	gmEmail = "LinkedUp Game <gm@linkedup.sfblockchainweek.io>"
)

// Client used to send emails
type Client interface {
	SendOnboardingEmail(*eb.AttendeeProfile, sdk.AccAddress, string, string) error
	SendRecoveryEmail(*eb.AttendeeProfile, int, string) error
}

type sesClient struct {
	ses         *ses.SES
	longyAppURL string
}

type mockClient struct {
	longyAppURL string
}

// NewMockClient creates a mock email client session wrapper that just logs
// the template parameters so that the application can run locally without
// actually sending email
func NewMockClient(longyAppURL string) (client Client, err error) {
	client = mockClient{
		longyAppURL: longyAppURL,
	}
	return
}

// NewSESClient creates a new email client session wrapper backed by Amazon SES
func NewSESClient(cfg client.ConfigProvider, localstack bool, longyAppURL string) (client Client, err error) {
	if localstack {
		client = sesClient{
			ses: ses.New(
				cfg,
				&aws.Config{
					Endpoint: aws.String("http://localstack:4579"),
				},
			),
			longyAppURL: longyAppURL,
		}
	} else {
		client = sesClient{
			ses:         ses.New(cfg),
			longyAppURL: longyAppURL,
		}
	}

	return
}

// SendOnboardingEmail will construct and send the email containing the initial
// onboarding message and URL with the given secret
func (c sesClient) SendOnboardingEmail(
	profile *eb.AttendeeProfile,
	attendeeAddr sdk.AccAddress,
	secret string,
	imageUploadURL string,
) error {
	redirectURI, err := makeOnboardingURI(c.longyAppURL, profile, attendeeAddr, secret, imageUploadURL)
	if err != nil {
		log.Errorf("unable to generate email URI: %s", err.Error())
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
	redirectURI, err := makeRecoveryURI(c.longyAppURL, id, token)

	if err != nil {
		return err
	}

	log.Tracef("sending recovery email to: %s", profile.Email)

	err = c.sendEmailWithURL(profile.Email, redirectURI, "linkedup-rekey")
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
		Source:       &gmEmail,
		Template:     &template,
		TemplateData: &templateData,
	})

	return
}

// SendOnboardingEmail will construct and send the email corresponding to onboarding the user
func (c mockClient) SendOnboardingEmail(
	profile *eb.AttendeeProfile,
	attendeeAddr sdk.AccAddress,
	secret string,
	imageUploadURL string,
) error {
	redirectURI, err := makeOnboardingURI(c.longyAppURL, profile, attendeeAddr, secret, imageUploadURL)
	if err != nil {
		return err
	}

	log.Infof("mock onboarding email with url: %s", redirectURI)
	return nil
}

func (c mockClient) SendRecoveryEmail(profile *eb.AttendeeProfile, id int, token string) error {
	redirectURI, err := makeRecoveryURI(c.longyAppURL, id, token)
	if err != nil {
		return err
	}

	log.Warnf("mock recovery email with url: %s", redirectURI)
	return nil
}

/** Helpers **/

func makeRecoveryURI(clientURL string, id int, token string) (string, error) {
	baseURL, err := url.Parse(fmt.Sprintf("%s/recover", clientURL))
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Add("id", strconv.Itoa(id))
	params.Add("token", token)

	baseURL.RawQuery = params.Encode()

	return baseURL.String(), nil
}

func makeOnboardingURI(
	clientURL string,
	profile *eb.AttendeeProfile,
	attendeeAddr sdk.AccAddress,
	secret string,
	imageUploadURL string,
) (string, error) {
	jsonProfileData, err := json.Marshal(profile)
	if err != nil {
		log.WithError(err).Error("attendee profile serialization")
		return "", err
	}

	baseURL, err := url.Parse(fmt.Sprintf("%s/claim", clientURL))

	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Add("attendee", attendeeAddr.String())
	params.Add("profile", base64.StdEncoding.EncodeToString(jsonProfileData))
	params.Add("secret", secret)
	params.Add("avatar", imageUploadURL)

	baseURL.RawQuery = params.Encode()

	return baseURL.String(), nil
}
