package eventbrite

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const (
	// EventEnvKey -
	EventEnvKey = "EVENTBRTIE_EVENT"
	// AuthEnvKey -
	AuthEnvKey = "EVENTBRITE_AUTH"

	urlFormat  = "https://www.eventbriteapi.com/v3/events/%d/attendees/?page=%d"
	authFormat = "Bearer %s"
)

var netClient = &http.Client{
	Timeout: 5 * time.Second,
}

// AttendeeProfile -
type AttendeeProfile struct {
	ID int `json:"-"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// GetAttendees -
func GetAttendees(eventID int, authToken string) ([]AttendeeProfile, error) {
	currentPage := 1
	hasMore := false
	var attendees []AttendeeProfile
	var profiles []AttendeeProfile // used for merging if there is more than 1 page

	var err error
	attendees, hasMore, err = fetchPage(eventID, authToken, currentPage)
	if err != nil {
		err = fmt.Errorf("page fetch from eventbrite: %s", err)
		return nil, err
	}
	currentPage++

	for hasMore {
		profiles, hasMore, err = fetchPage(eventID, authToken, currentPage)
		if err != nil {
			err = fmt.Errorf("page fetch from eventbrite: %s", err)
			return nil, err
		}

		currentPage++
		attendees = append(attendees, profiles...)
	}

	return attendees, nil
}

// GetAttendeesFromEnv -
func GetAttendeesFromEnv() ([]AttendeeProfile, error) {
	eventStr := os.Getenv(EventEnvKey)
	authToken := os.Getenv(AuthEnvKey)
	if len(eventStr) == 0 || len(authToken) == 0 {
		err := fmt.Errorf("%s and %s environment variables must be set to communicate with eventbrite",
			EventEnvKey, AuthEnvKey)
		return nil, err
	}

	eventID, err := strconv.Atoi(eventStr)
	if err != nil {
		err = fmt.Errorf("event id must be a positive number in decimal format: %s", err)
		return nil, err
	}

	return GetAttendees(eventID, authToken)
}

func fetchPage(eventID int, authToken string, page int) (attendees []AttendeeProfile, hasMore bool, err error) {
	type pageInfo struct {
		// only things we care about within the page struct
		HasMore bool `json:"has_more_items"`
	}
	type bodyResp struct {
		Page      pageInfo          `json:"pagination"`
		Attendees []json.RawMessage `json:"attendees"`
	}

	auth := fmt.Sprintf(authFormat, authToken)
	url, err := url.Parse(fmt.Sprintf(urlFormat, eventID, page))
	if err != nil {
		err = fmt.Errorf("parsing url: %s", err)
		return nil, false, err
	}
	req := &http.Request{
		URL:    url,
		Method: "GET",
		Header: map[string][]string{
			"Authorization": {auth},
		},
	}
	resp, err := netClient.Do(req)
	if err != nil {
		err = fmt.Errorf("eventbrite request delivery: %s", err)
		return nil, false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad response. status code=%d", resp.StatusCode)
		return nil, false, err
	}

	// decode body
	decoder := json.NewDecoder(resp.Body)
	var data bodyResp
	if err := decoder.Decode(&data); err != nil {
		err = fmt.Errorf("reading request body: %s", err)
		return nil, false, err
	}

	// retrieve attendees
	numAttendees := len(data.Attendees)
	attendees = make([]AttendeeProfile, numAttendees)
	for i := 0; i < numAttendees; i++ {
		attendee, err := getProfile(data.Attendees[i])
		if err != nil {
			return nil, false, err
		}

		attendees[i] = *attendee
	}

	hasMore = data.Page.HasMore
	err = nil

	return attendees, hasMore, err
}

func getProfile(body json.RawMessage) (*AttendeeProfile, error) {
	/** Parse the entire attendee body **/
	var jsonResp map[string]json.RawMessage
	if err := json.Unmarshal(body, &jsonResp); err != nil {
		err = fmt.Errorf("parsing eventbrite attendee: %s", err)
		return nil, err
	}

	var idStr string
	if err := json.Unmarshal(jsonResp["id"], &idStr); err != nil {
		return nil, fmt.Errorf("unable to parse attendee id: %s", err)
	}

	attendeeID, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, fmt.Errorf("expected attendee id to be a number")
	}

	/** Extract specifically the profile key of the response **/
	var jsonProfile map[string]json.RawMessage
	profileData, ok := jsonResp["profile"]
	if !ok {
		return nil, fmt.Errorf("attendee missing profile")
	} else if err := json.Unmarshal(profileData, &jsonProfile); err != nil {
		err = fmt.Errorf("parsing attendee profile: %s", err)
		return nil, err
	}

	/** Decode the struct into the fields we want **/
	var profile AttendeeProfile
	err = json.Unmarshal(profileData, &profile)
	if err != nil {
		err = fmt.Errorf("decoding profile: %s", err)
		return nil, err
	}
	profile.ID = attendeeID

	return &profile, nil
}
