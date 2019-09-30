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
	attendeePrefix = []byte{0x0}
)

func AttendeeKeyByID(id string) []byte {
	return prefixKey(attendeePrefix, []byte(id))
}

func prefixKey(prefix, key []byte) []byte {
	buf := new(bytes.Buffer)
	buf.Write(prefix)
	buf.WriteString("::")
	buf.Write(key)

	return buf.Bytes()
}
