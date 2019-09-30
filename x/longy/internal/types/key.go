package types

import (
	"bytes"
)

const (
	// ModuleName is the name of this module
	ModuleName = "longy"

	// RouterKey is the package route
	RouterKey = ModuleName
)

var (
	attendeePrefix  = []byte{0x0}
	masterKeyPrefix = []byte{0x1}
)

// AttendeeKeyByID will construct the appropriate key for the attendee with `id`
func AttendeeKeyByID(id string) []byte {
	return prefixKey(attendeePrefix, []byte(id))
}

// MasterKey will return the store key for the master key
func MasterKey() []byte {
	return prefixKey(masterKeyPrefix, []byte("master"))
}

func prefixKey(prefix, key []byte) []byte {
	buf := new(bytes.Buffer)
	buf.Write(prefix)
	buf.WriteString("::")
	buf.Write(key)

	return buf.Bytes()
}
