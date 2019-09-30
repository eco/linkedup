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
	attendeePrefix  = []byte{0x0}
	masterKeyPrefix = []byte{0x1}
	rewardKeyPrefix = []byte{0x2}
)

// AttendeeKeyByID will construct the appropriate key for the attendee with `id`
func AttendeeKey(addr sdk.AccAddress) []byte {
	return prefixKey(attendeePrefix, addr[:])
}

// MasterKey will return the store key for the master key
func MasterKey() []byte {
	return masterKeyPrefix
}

// RewardKey creates the key for reward of type 'kind`
func RewardKey(kind string) []byte {
	return prefixKey(rewardKeyPrefix, []byte(kind))
}

func prefixKey(prefix, key []byte) []byte {
	buf := new(bytes.Buffer)
	buf.Write(prefix)
	buf.WriteString("::")
	buf.Write(key)

	return buf.Bytes()
}
