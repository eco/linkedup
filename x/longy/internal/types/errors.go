package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	//LongyCodeSpace is the codespace  type for errors
	LongyCodeSpace sdk.CodespaceType = ModuleName

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
	//QRCodeInvalid is the code when a scan message has an invalid QR code, ie not a positive integer
	QRCodeInvalid
	//AttendeeUnclaimed is the code when
	AttendeeUnclaimed
	//AttendeeAlreadyClaimed is the code when
	AttendeeAlreadyClaimed
	//InvalidCommitmentReveal is the code when
	InvalidCommitmentReveal
	//ScanAccountsSame is the code when a scan of with the same 1 account is attempted
	ScanAccountsSame
	//AccountAddressEmpty is the code when an AccAddress is the empty address
	AccountAddressEmpty

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

// ErrQRCodeInvalid occurs when a scan message has an invalid QR code, ie not a positive integer
func ErrQRCodeInvalid(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, QRCodeInvalid, format, args...)
}

// ErrAttendeeUnclaimed indicates a attendee that is unclaimed
func ErrAttendeeUnclaimed() sdk.Error {
	return sdk.NewError(LongyCodeSpace, AttendeeUnclaimed, "attendee unclaimed")
}

// ErrAttendeeAlreadyClaimed indicates a attendee that is unclaimed
func ErrAttendeeAlreadyClaimed() sdk.Error {
	return sdk.NewError(LongyCodeSpace, AttendeeAlreadyClaimed, "attendee claimed")
}

// ErrInvalidCommitmentReveal indicates that the reveal is incorrect for the commitment
func ErrInvalidCommitmentReveal() sdk.Error {
	return sdk.NewError(LongyCodeSpace, InvalidCommitmentReveal, "reveal to the commitment is incorrect")
}

//ErrScanAccountsSame occurs when a scan of with the same 1 account is attempted
func ErrScanAccountsSame(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, ScanAccountsSame, format, args...)
}

//ErrAccountAddressEmpty occurs when an AccAddress is the empty address
func ErrAccountAddressEmpty(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, AccountAddressEmpty, format, args...)
}

//ErrDefault occurs when a random error occurs that we do not provide a unique code to
func ErrDefault(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, DefaultError, format, args...)
}
