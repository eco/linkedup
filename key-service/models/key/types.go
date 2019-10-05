package key

// StoredKey represents a record associating some key data with an email
// address in the application database.
type StoredKey struct {
	Email    string
	KeyData  []byte
}
