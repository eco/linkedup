package eventbrite

import (
	"github.com/eco/longy/eventbrite"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	log = logrus.WithField("module", "eventbrite-session")
)

// Session to hook into the eventbrite API
type Session struct {
	sync.Mutex
	ticker *time.Ticker

	numAttendees int
	attendees    map[int]eventbrite.AttendeeProfile
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

	ticker := time.NewTicker(5 * time.Minute)
	log.Info("eventbrite polling setup for 5 minute intervals")

	session := &Session{
		Mutex:  sync.Mutex{},
		ticker: ticker,

		numAttendees: len(attendees),
		attendees:    make(map[int]eventbrite.AttendeeProfile),
	}

	go session.poll(ticker, eventID, authToken)

	session.Lock()
	for _, a := range attendees {
		session.attendees[a.ID] = a
	}
	session.Unlock()

	return session, nil
}

// Close will release resources
func (s *Session) Close() {
	log.Info("ending polling with eventbrite")
	s.ticker.Stop()
}

// AttendeeProfile -
func (s *Session) AttendeeProfile(id int) (*eventbrite.AttendeeProfile, bool) {
	profile, ok := s.attendees[id]
	return &profile, ok
}

func (s *Session) poll(ticker *time.Ticker, eventID int, authToken string) {
	for range ticker.C {
		attendees, err := eventbrite.GetAttendees(eventID, authToken)
		if err != nil {
			log.WithError(err).Warn("error polling eventbrite")
			continue
		}

		if len(attendees) != s.numAttendees {
			// there are updates
			newMap := make(map[int]eventbrite.AttendeeProfile)
			for _, a := range attendees {
				newMap[a.ID] = a
			}

			s.Lock()
			s.attendees = newMap
			s.numAttendees = len(attendees)
			s.Unlock()

			log.Info("updated cached attendee list")
		}

		// no updates, continue
	}
}
