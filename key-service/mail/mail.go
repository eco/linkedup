package mail

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/ses"
	sdk "github.com/cosmos/cosmos-sdk/types"
	eb "github.com/eco/longy/eventbrite"
	"github.com/eco/longy/key-service/models"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("module", "mail")

//GmEmail is the authorized from email for AWS SES email sending
var GmEmail = "LinkedUp Game <gm@linkedup.sfblockchainweek.io>"

//var GmEmail = "LinkedUp Game <linkedup@sfbw.io>"

const (
	onboardingRedirectTemplate = "linkedup-onboarding"
	rekeyRedirectTemplate      = "linkedup-rekey"
	verificationTemplate       = "linkedup-verification"
	retrieveTemplate           = "linkedup-retrieve"
)

// Client used to send emails
type Client interface {
	SendOnboardingEmail(*models.DatabaseContext, sdk.AccAddress, *eb.AttendeeProfile, string, string) error
	SendRecoveryEmail(*models.DatabaseContext, *eb.AttendeeProfile, string, string) error

	SendVerificationEmail(*models.DatabaseContext, string, string) error

	SendExportEmail(db *models.DatabaseContext, attendeeEmail string, id int, token string) error

	SendAttendeeSharedInfoEmail(ctx *models.DatabaseContext, attendeeEmail string, sharedInfo string) error
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
	db *models.DatabaseContext,
	attendeeAddr sdk.AccAddress,
	profile *eb.AttendeeProfile,
	secret string,
	imageUploadURL string,
) error {
	redirectURI, err := makeOnboardingURI(c.longyAppURL, attendeeAddr, profile, secret, imageUploadURL)
	if err != nil {
		log.Errorf("unable to generate email URI: %s", err.Error())
		return err
	}

	log.Tracef("sending onboarding email to: %s", profile.Email)

	err = c.sendEmailWithURL(db, profile.Email, redirectURI, onboardingRedirectTemplate)
	if err != nil {
		log.WithError(err).Errorf("unable to send onboarding email to %s", profile.Email)
	}

	return err
}

// SendRecoveryEmail will construct and send the email containing the account
// recovery message and URL with the given secret
func (c sesClient) SendRecoveryEmail(
	db *models.DatabaseContext,
	profile *eb.AttendeeProfile,
	id string,
	token string,
) error {
	redirectURI, err := makeRecoveryURI(c.longyAppURL, id, token)
	if err != nil {
		return err
	}

	log.Tracef("sending recovery email to: %s", profile.Email)

	err = c.sendEmailWithURL(db, profile.Email, redirectURI, rekeyRedirectTemplate)
	if err != nil {
		log.WithError(err).Errorf("unable to send recovery email to %s", profile.Email)
	}

	return err
}

func (c sesClient) SendVerificationEmail(
	db *models.DatabaseContext,
	dest string,
	token string,
) error {
	if db.GetBlacklistEntry(dest) {
		log.WithField("dest", dest).Trace("refusing to email to blacklisted address")
		return nil
	}
	var template = verificationTemplate
	templateData := fmt.Sprintf("{\"token\":\"%s\"}", token)
	_, err := c.ses.SendTemplatedEmail(&ses.SendTemplatedEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{&dest},
		},
		Source:       &GmEmail,
		Template:     &template,
		TemplateData: &templateData,
	})

	return err
}

func (c sesClient) SendExportEmail(db *models.DatabaseContext, dstEmail string, id int, token string) error {
	//https://linkedup.sfbw.io/s/export/index.html?id=1284763463&token=584353
	//EmailExportUrlBase is the base url for info export
	link := fmt.Sprintf("%s/s/export/index.html?id=%d&token=%s", c.longyAppURL, id, token)

	err := c.sendEmailWithURL(db, dstEmail, link, retrieveTemplate)
	if err != nil {
		log.WithError(err).Errorf("unable to send recovery email to %s", dstEmail)
	}

	return err
}

//SendAttendeeSharedInfoEmail sends the shared info to an attendee
func (c sesClient) SendAttendeeSharedInfoEmail(
	db *models.DatabaseContext,
	attendeeEmail string,
	sharedInfo string) error {
	if db.GetBlacklistEntry(attendeeEmail) {
		log.WithField("dest", attendeeEmail).Trace("refusing to email to blacklisted address")
		return nil
	}

	_, err := sendRaw(c.ses, attendeeEmail, sharedInfo)

	return err
}

func (c sesClient) sendEmailWithURL(
	db *models.DatabaseContext,
	dest string,
	url string,
	template string,
) error {
	if db.GetBlacklistEntry(dest) {
		log.WithField("dest", dest).Trace("refusing to email to blacklisted address")
		return nil
	}
	templateData := fmt.Sprintf("{\"url\":\"%s\"}", url)
	_, err := c.ses.SendTemplatedEmail(&ses.SendTemplatedEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{&dest},
		},
		Source:       &GmEmail,
		Template:     &template,
		TemplateData: &templateData,
	})

	return err
}

// SendOnboardingEmail will construct and send the email corresponding to onboarding the user
func (c mockClient) SendOnboardingEmail(
	db *models.DatabaseContext,
	attendeeAddr sdk.AccAddress,
	profile *eb.AttendeeProfile,
	secret string,
	imageUploadURL string,
) error {
	redirectURI, err := makeOnboardingURI(c.longyAppURL, attendeeAddr, profile, secret, imageUploadURL)
	if err != nil {
		return err
	}

	log.Infof("mock onboarding email with url: %s", redirectURI)
	return nil
}

func (c mockClient) SendRecoveryEmail(
	db *models.DatabaseContext,
	profile *eb.AttendeeProfile,
	id string,
	token string,
) error {
	redirectURI, err := makeRecoveryURI(c.longyAppURL, id, token)
	if err != nil {
		return err
	}

	log.Warnf("mock recovery email with url: %s", redirectURI)
	return nil
}

func (c mockClient) SendVerificationEmail(
	db *models.DatabaseContext,
	dest string,
	token string,
) error {
	log.Warnf("mock verification token: %s", token)
	return nil
}

func (c mockClient) SendExportEmail(db *models.DatabaseContext, dstEmail string, id int, token string) error {
	link := fmt.Sprintf("%s/s/export/index.html?id=%d&token=%s", c.longyAppURL, id, token)
	log.Warnf("mock SendExportEmail : %s", link)
	return nil
}

func (c mockClient) SendAttendeeSharedInfoEmail(
	ctx *models.DatabaseContext,
	attendeeEmail string,
	sharedInfo string) error {
	log.Warnf("mock attendee share info : %s", sharedInfo)
	return nil
}

/** Helpers **/

func makeRecoveryURI(clientURL string, id string, token string) (string, error) {
	baseURL, err := url.Parse(fmt.Sprintf("%s/recover", clientURL))
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Add("id", id)
	params.Add("token", token)

	baseURL.RawQuery = params.Encode()

	return baseURL.String(), nil
}

func makeOnboardingURI(
	clientURL string,
	attendeeAddr sdk.AccAddress,
	profile *eb.AttendeeProfile,
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
