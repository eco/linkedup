package eventbrite

import (
	"github.com/eco/longy/eventbrite"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithField("module", "eventbrite-session")
)

// Session to hook into the eventbrite API
type Session struct {
	attendees map[int]eventbrite.AttendeeProfile
}

// CreateSession to interact iwht the EventBrite Event APIs. The constructed
// session has a default timeout of 10 seconds
func CreateSession(eventID int, authToken string) (*Session, error) {
	log.WithField("auth", authToken).
		WithField("event", eventID).
		Info("eventbrite session created")

	attendees, err := eventbrite.GetAttendees(eventID, authToken)
	if err != nil {
		return nil, err
	}

	log.Infof("retrieved %d attendees from eventbrite", len(attendees))

	session := &Session{
		attendees: make(map[int]eventbrite.AttendeeProfile),
	}

	for _, a := range attendees {
		session.attendees[a.ID] = a
	}

	return session, nil
}

func (s *Session) AttendeeProfile(id int) (*eventbrite.AttendeeProfile, bool) {
	profile, ok := s.attendees[id]
	return &profile, ok
}
