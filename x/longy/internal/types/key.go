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
	// AttendeePrefix is the prefix for the attendee type
	AttendeePrefix = []byte{0x0}
	// ScanPrefix is the prefix for the scan type
	ScanPrefix = []byte{0x1}
	// PrizePrefix is the prefix for the prize type
	PrizePrefix = []byte{0x2}
	// MasterKeyPrefix is the prefix for storing the address of the service account
	MasterKeyPrefix = []byte{0x3}
	// BonusKeyPrefix is the prefix for retrieving the active bonus
	BonusKeyPrefix = []byte{0x4}

	keySeparator = []byte("::")
)

// AttendeeKey will construct the appropriate key for the attendee with `id`
func AttendeeKey(addr sdk.AccAddress) []byte {
	return prefixKey(AttendeePrefix, addr[:])
}

// ScanKey returns the prefixed key for managing scans in the store
func ScanKey(id []byte) []byte {
	return prefixKey(ScanPrefix, id)
}

// PrizeKey returns the prefixed key for managing prizes in the store
func PrizeKey(id []byte) []byte {
	return prefixKey(PrizePrefix, id)
}

// MasterKey will return the store key for the master key
func MasterKey() []byte {
	return MasterKeyPrefix
}

// BonusKey -
func BonusKey() []byte {
	return BonusKeyPrefix
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
	return append(pre, keySeparator...)
}
