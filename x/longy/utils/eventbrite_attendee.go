package utils

import (
	"github.com/eco/longy/util"
	"github.com/eco/longy/x/longy"
	"strings"
)

const (
	//TicketSponsorNameLowerCase is the sponsors ticket type
	TicketSponsorNameLowerCase = "sponsors"
	//TicketSpeakerCescNameLowerCase is the cesc speakers ticket type
	TicketSpeakerCescNameLowerCase = "cesc speakers"
	//TicketSpeakerEpicenterNameLowerCase is the epicenter speakers ticket type
	TicketSpeakerEpicenterNameLowerCase = "epicenter speakers"
)

//EventbriteAttendees is the array of attendees that the api returns for processing
type EventbriteAttendees []EventbriteAttendee

// EventbriteAttendee is the attendee structure in the genesis file
type EventbriteAttendee struct {
	ID              string            `json:"id"`
	TicketClassName string            `json:"ticket_class_name"`
	Profile         EventbriteProfile `json:"profile"` //gets the full info of the account
}

//ToGenesisAttendee turns the eventbrite type to our local type
func (e *EventbriteAttendee) ToGenesisAttendee() longy.Attendee {
	return longy.Attendee{
		ID:      e.ID,
		Address: util.IDToAddress(e.ID),
		Sponsor: e.IsSponsorTicket(),
		Name:    e.Profile.Name,
		// can also add Profile to the Attendee's struc to show all
	}
}

//IsSponsorTicket checks to see if the ticket type is of a speaker or sponsor that gets special point bonuses
func (e *EventbriteAttendee) IsSponsorTicket() bool {
	switch strings.ToLower(e.TicketClassName) {
	case TicketSponsorNameLowerCase:
		fallthrough
	case TicketSpeakerCescNameLowerCase:
		fallthrough
	case TicketSpeakerEpicenterNameLowerCase:
		return true
	}
	return false
}

// EventbriteProfile is the profile of the attendee from eventbrite
type EventbriteProfile struct {
	Name     string `json:"name"`
	Company  string `json:"company"`
	Email    string `json:"email"`
	JobTitle string `json:"job_title"`
}
