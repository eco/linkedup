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

// AttendeeProfile -
type AttendeeProfile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// AttendeeProfile retrieves the email associated with the eventbrite account for `id`
func (s Session) AttendeeProfile(id string) (*AttendeeProfile, error) {
	/** Construct the appropriate url **/
	host := "https://www.eventbriteapi.com"
	path := fmt.Sprintf("/v3/events/%d/attendees/%s/", s.eventID, id)
	auth := fmt.Sprintf("Bearer %s", s.authToken)

	url, err := url.Parse(host + path)
	if err != nil {
		log.Warnf("bad event url: %s", host+path)
		return nil, ErrInternal
	}

	/** Create the request to issue **/
	req := &http.Request{
		URL:    url,
		Method: "GET",
		Header: map[string][]string{
			"Authorization": {auth},
		},
	}
	resp, err := s.netClient.Do(req)
	if err != nil {
		log.WithError(err).Error("eventbrite api request delivery")
		return nil, ErrInternal
	}
	defer resp.Body.Close()

	/** Read the EventBrite response **/
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
		log.WithField("status_code", resp.StatusCode).WithField("attendee_id", id).
			Error("bad eventbrite api response")

		return nil, ErrInternal
	}

	return getProfileFromBody(resp.Body)
}

func getProfileFromBody(body io.Reader) (*AttendeeProfile, error) {
	/** Parse the raw request **/
	var jsonResp map[string]json.RawMessage
	d := json.NewDecoder(body)
	if err := d.Decode(&jsonResp); err != nil {
		log.WithError(err).Error("parsing eventbrite response")
		return nil, ErrInternal
	}

	/** Extract specifically the profile key of the response **/
	var jsonProfile map[string]json.RawMessage
	profileData, ok := jsonResp["profile"]
	if !ok {
		log.Error("eventbrite response missing attendee profile")
		return nil, ErrInternal
	} else if err := json.Unmarshal(profileData, &jsonProfile); err != nil {
		log.WithError(err).Error("parsing attendee profile")
		return nil, ErrInternal
	}

	/** Decode the struct into the fields we want **/
	var profile AttendeeProfile
	err := json.Unmarshal(profileData, &profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}
