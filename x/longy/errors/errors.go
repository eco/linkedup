package errors

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/types"
)

const (
	//LongyCodeSpace is the codespace  type for errors
	LongyCodeSpace sdk.CodespaceType = types.ModuleName

	//ItemNotFound is the code for no item
	ItemNotFound sdk.CodeType = iota + 1
	//AttendeeNotFound is the code for when the attendee cannot be found in the keeper
	AttendeeNotFound
	//InsufficientPrivileges is the code for when a transaction signer doesn't have the necessary privilege
	InsufficientPrivileges
	//GenesisServiceAddressEmpty is the code for when the service account address is not set in the genesis file
	GenesisServiceAddressEmpty
	//GenesisAttendeesEmpty is the code for when the attendees are not set in the genesis file
	GenesisAttendeesEmpty
	//EventbriteEnvVariableNotSet is the code for when the eventbrite environmental var containing the auth key is
	//not set
	EventbriteEnvVariableNotSet
	//NetworkResponseError is the code for any network response that is not what is expected, ie 200/201
	NetworkResponseError
	//AttendeeCountMismatch is the code for when there is a mis match between the expected and received number of
	//attendees from the indexing of the eventbrite calls
	AttendeeCountMismatch
	//GenesisServiceAccountInvalid is the code when the service account bech32 address is invalidly passed to gen
	GenesisServiceAccountInvalid
	//GenesisServiceAccountNotPresent is the code when the service account is not found in the genesis accounts
	GenesisServiceAccountNotPresent

	//DefaultError is the code for when a random error occurs that we do not provide a unique code to
	DefaultError
)

//ErrItemNotFound occurs when we cannot find an item in the default store
func ErrItemNotFound(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, ItemNotFound, format, args...)
}

//ErrAttendeeNotFound occurs when we cannot find the attendee
func ErrAttendeeNotFound(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, AttendeeNotFound, format, args...)
}

//ErrInsufficientPrivileges occurs when we cannot find the attendee
func ErrInsufficientPrivileges(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, InsufficientPrivileges, format, args...)
}

//ErrGenesisServiceAddressEmpty occurs when the re-key service address is not set in the genesis file
func ErrGenesisServiceAddressEmpty(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, GenesisServiceAddressEmpty, format, args...)
}

//ErrGenesisAttendeesEmpty occurs when the attendees are not set in the genesis file
func ErrGenesisAttendeesEmpty(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, GenesisAttendeesEmpty, format, args...)
}

//ErrEventbriteEnvVariableNotSet occurs when the attendees are not set in the genesis file
func ErrEventbriteEnvVariableNotSet(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, EventbriteEnvVariableNotSet, format, args...)
}

//ErrNetworkResponseError occurs when network response that is not what is expected, ie 200/201
func ErrNetworkResponseError(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, NetworkResponseError, format, args...)
}

//ErrAttendeeCountMismatch occurs when network response that is not what is expected, ie 200/201
func ErrAttendeeCountMismatch(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, AttendeeCountMismatch, format, args...)
}

//ErrGenesisServiceAccountInvalid occurs when the service account bech32 address is invalidly passed to gen
func ErrGenesisServiceAccountInvalid(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, GenesisServiceAccountInvalid, format, args...)
}

//ErrGenesisServiceAccountNotPresent occurs when the service account is not found in the genesis accounts
func ErrGenesisServiceAccountNotPresent(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, GenesisServiceAccountNotPresent, format, args...)
}

//ErrDefault occurs when a random error occurs that we do not provide a unique code to
func ErrDefault(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, DefaultError, format, args...)
}
