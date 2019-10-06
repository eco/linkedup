package eventbrite

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"time"
)

var (
	// internal
	log = logrus.WithField("module", "eventbrite")

	// ErrInternal -
	ErrInternal = errors.New("internal error")
)

// Session to hook into the eventbrite API
type Session struct {
	authToken string
	eventID   int
	netClient http.Client
}

// CreateSession to interact iwht the EventBrite Event APIs. The constructed
// session has a default timeout of 10 seconds
func CreateSession(authToken string, eventID int) Session {
	log.WithField("auth", authToken).
		WithField("event", eventID).
		Info("eventbrite session created")

	return Session{
		authToken: authToken,
		eventID:   eventID,
		netClient: http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

// WithTimeout returns a new `Session` with the corresponding `time` timeout
func (s Session) WithTimeout(time time.Duration) Session {
	s.netClient = http.Client{
		Timeout: time,
	}

	return s
}

// EmailFromAttendeeID retrieves the email associated with the eventbrite account for `id`
func (s Session) EmailFromAttendeeID(id string) (string, error) {
	host := "https://eventbrite.com"
	path := fmt.Sprintf("/v3/%d/attendees/%s/", s.eventID, id)
	auth := fmt.Sprintf("Bearer %s", s.authToken)

	url, err := url.Parse(host + path)
	if err != nil {
		log.Warnf("bad event url: %s", host+path)
		return "", ErrInternal
	}

	// create the http request
	req := &http.Request{
		URL:    url,
		Method: "GET",
		Header: map[string][]string{
			"Authorization": {auth},
		},
	}

	resp, err := s.netClient.Do(req)
	if err != nil {
		log.WithError(err).Warn("eventbrite api request delivery")
		return "", ErrInternal
	}
	defer resp.Body.Close() //nolint
	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	} else if resp.StatusCode != http.StatusOK {
		// TODO: emperically check the different response types. id does not exist a 403 (NotFound)?
		// Log a warning?
		log.Warnf("bad eventbrite api response, code=%d, attendee_id=%s", resp.StatusCode, id)
		return "", ErrInternal
	}

	return getEmailFromBody(resp.Body)
}

func getEmailFromBody(body io.Reader) (string, error) {
	var jsonResp map[string]json.RawMessage
	d := json.NewDecoder(body)
	if err := d.Decode(&jsonResp); err != nil {
		log.WithError(err).Error("parsing eventbrite response")
		return "", ErrInternal
	}

	var jsonProfile map[string]json.RawMessage
	profileData, ok := jsonResp["profile"]
	if !ok {
		log.Error("eventbrite response missing attendee profile")
		return "", ErrInternal
	} else if err := json.Unmarshal(profileData, &jsonProfile); err != nil {
		log.WithError(err).Error("parsing attendee profile")
		return "", ErrInternal
	}

	var email string
	emailData, ok := jsonProfile["email"]
	if !ok {
		log.Error("attendee profile missing email")
		return "", ErrInternal
	} else if err := json.Unmarshal(emailData, &email); err != nil {
		log.WithError(err).Error("parsing email")
		return "", ErrInternal
	}

	return email, nil
}
