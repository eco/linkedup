package types

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of this module
	ModuleName = "longy"

	// StoreKey is the key used to access the store
	StoreKey = ModuleName

	// RouterKey is the package route
	RouterKey = ModuleName
)

var (
	//AttendeePrefix is the prefix for the attendee type
	AttendeePrefix = []byte{0x0}
	//ScanPrefix is the prefix for the scan type
	ScanPrefix = []byte{0x1}
	//InfoPrefix is the prefix for the info type
	InfoPrefix = []byte{0x2}
	//MasterKeyPrefix is the prefix for storing the public address of the service account
	MasterKeyPrefix = []byte{0x3}
	//KeySeparator is the separator between the prefix and the type key
	KeySeparator = []byte("::")
)

// AttendeeKey will construct the appropriate key for the attendee with `id`
func AttendeeKey(addr sdk.AccAddress) []byte {
	return prefixKey(AttendeePrefix, addr[:])
}

//ScanKey returns the prefixed key for managing scans in the store
func ScanKey(id []byte) []byte {
	return prefixKey(ScanPrefix, id)
}

//InfoKey returns the prefixed key for managing info in the store
func InfoKey(id []byte) []byte {
	return prefixKey(InfoPrefix, id)
}

// MasterKey will return the store key for the master key
func MasterKey() []byte {
	return MasterKeyPrefix
}

//nolint:gosec
func prefixKey(pre []byte, key []byte) []byte {
	buf := new(bytes.Buffer)
	buf.Write(Prefix(pre))
	buf.Write(key)

	return buf.Bytes()
}

//Prefix returns the prefix for a given pre key
func Prefix(pre []byte) []byte {
	return append(pre, KeySeparator...)
}
