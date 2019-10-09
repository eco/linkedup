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
