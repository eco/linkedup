package models

// StoredKey represents a record associating some attendee information
type storedInfo struct {
	ID   int
	Data []byte
}

// StoredAuth represents a record associating an auth token
type storedAuth struct {
	ID        int
	AuthToken string
}

// email is the id <-> email override that we use for attendees who's eventbrite emails are not set correctly
type storeEmail struct {
	ID    int
	Email string
}

// provides a way to manage an email blacklist for bounced emails
type blacklistEmail struct {
	Email       string
	Blacklisted bool
}
