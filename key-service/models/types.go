package models

// StoredKey represents a record associating some key data with an email
// address in the application database.
type storedKey struct {
	Email   string
	KeyData []byte
}
