package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	//LongyCodeSpace is the codespace  type for errors
	LongyCodeSpace sdk.CodespaceType = ModuleName

	// ItemNotFound is the code for no item
	ItemNotFound sdk.CodeType = iota + 1

	// AttendeeNotFound is the code for when the attendee cannot be found in the keeper
	AttendeeNotFound

	// ScanNotFound is the code when we cannot find a scan in the keeper with the given id
	ScanNotFound

	//InfoNotFound is the code when we cannot find an info in the keeper with the given id
	InfoNotFound

	// InsufficientPrivileges is the code for when a transaction signer doesn't have the necessary privilege
	InsufficientPrivileges

	// GenesisServiceAddressEmpty is the code for when the service account address is not set in the genesis file
	GenesisServiceAddressEmpty

	// GenesisAttendeesEmpty is the code for when the attendees are not set in the genesis file
	GenesisAttendeesEmpty

	// EventbriteEnvVariableNotSet is the code for when the eventbrite environmental var containing the auth key is
	//not set
	EventbriteEnvVariableNotSet

	// NetworkResponseError is the code for any network response that is not what is expected, ie 200/201
	NetworkResponseError

	// AttendeeCountMismatch is the code for when there is a mis match between the expected and received number of
	// attendees from the indexing of the eventbrite calls
	AttendeeCountMismatch

	// GenesisServiceAccountInvalid is the code when the service account bech32 address is invalidly passed to gen
	GenesisServiceAccountInvalid

	// GenesisServiceAccountNotPresent is the code when the service account is not found in the genesis accounts
	GenesisServiceAccountNotPresent

	// QRCodeInvalid is the code when
	QRCodeInvalid

	// AttendeeClaimed is the code when
	AttendeeClaimed

	// AttendeeKeyed is the code when the attendee has already been key'd by the rekey service
	AttendeeKeyed

	//InvalidCommitmentReveal is the code when
	InvalidCommitmentReveal

	// AccountsSame is the code when a scan of with the same 1 account is attempted
	AccountsSame

	// AccountAddressEmpty is the code when an AccAddress is the empty address
	AccountAddressEmpty

	// ScanQRAlreadyOccurred is the code for when the scan message has already been sent by the scanner or the scan
	// is complete for those two parties
	ScanQRAlreadyOccurred
	//ScanNotComplete is the code for when a scan is not complete
	ScanNotComplete
	//DataCannotBeEmpty is the code for when the info data in a message is empty
	DataCannotBeEmpty
	//DataSizeOverLimit is the code for when the info data is above the size limit
	DataSizeOverLimit
	//CantShareWithSelf is the code for when an attendee tries to share info with themselves
	CantShareWithSelf
	//InfoAlreadyExists is the code for when someone tries to share info with a person more than once
	InfoAlreadyExists
	//InvalidShareForScan is the code for when someone tries to share info when the corresponding scan is not complete
	InvalidShareForScan

	// DefaultError is the code for when a random error occurs that we do not provide a unique code to
	DefaultError
)

// ErrItemNotFound occurs when we cannot find an item in the default store
func ErrItemNotFound(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, ItemNotFound, format, args...)
}

// ErrScanNotFound occurs when we cannot find a scan in the keeper with the given id
func ErrScanNotFound(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, ScanNotFound, format, args...)
}

// ErrInfoNotFound occurs when we cannot find an info in the keeper with the given id
func ErrInfoNotFound(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, InfoNotFound, format, args...)
}

// ErrAttendeeNotFound occurs when we cannot find the attendee
func ErrAttendeeNotFound(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, AttendeeNotFound, format, args...)
}

// ErrInsufficientPrivileges occurs when we cannot find the attendee
func ErrInsufficientPrivileges(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, InsufficientPrivileges, format, args...)
}

// ErrGenesisServiceAddressEmpty occurs when the re-key service address is not set in the genesis file
func ErrGenesisServiceAddressEmpty(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, GenesisServiceAddressEmpty, format, args...)
}

// ErrGenesisAttendeesEmpty occurs when the attendees are not set in the genesis file
func ErrGenesisAttendeesEmpty(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, GenesisAttendeesEmpty, format, args...)
}

// ErrEventbriteEnvVariableNotSet occurs when the attendees are not set in the genesis file
func ErrEventbriteEnvVariableNotSet(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, EventbriteEnvVariableNotSet, format, args...)
}

// ErrNetworkResponseError occurs when network response that is not what is expected, ie 200/201
func ErrNetworkResponseError(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, NetworkResponseError, format, args...)
}

// ErrAttendeeCountMismatch occurs when network response that is not what is expected, ie 200/201
func ErrAttendeeCountMismatch(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, AttendeeCountMismatch, format, args...)
}

// ErrGenesisServiceAccountInvalid occurs when the service account bech32 address is invalidly passed to gen
func ErrGenesisServiceAccountInvalid(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, GenesisServiceAccountInvalid, format, args...)
}

// ErrGenesisServiceAccountNotPresent occurs when the service account is not found in the genesis accounts
func ErrGenesisServiceAccountNotPresent(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, GenesisServiceAccountNotPresent, format, args...)
}

// ErrQRCodeInvalid occurs when a scan message has an invalid QR code, ie not a positive integer
func ErrQRCodeInvalid(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, QRCodeInvalid, format, args...)
}

// ErrAttendeeClaimed indicates a attendee that is unclaimed
func ErrAttendeeClaimed(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, AttendeeClaimed, format, args...)
}

// ErrAttendeeKeyed indicated the attendee has already keyed their account
func ErrAttendeeKeyed(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, AttendeeKeyed, format, args...)
}

// ErrInvalidCommitmentReveal indicates that the reveal is incorrect for the commitment
func ErrInvalidCommitmentReveal(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, InvalidCommitmentReveal, format, args...)
}

// ErrAccountsSame occurs when a scan of with the same 1 account is attempted
func ErrAccountsSame(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, AccountsSame, format, args...)
}

// ErrAccountAddressEmpty occurs when an AccAddress is the empty address
func ErrAccountAddressEmpty(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, AccountAddressEmpty, format, args...)
}

// ErrScanQRAlreadyOccurred occurs when the scan message has already been sent by the scanner or the scan
//is complete for those two parties
func ErrScanQRAlreadyOccurred(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, ScanQRAlreadyOccurred, format, args...)
}

//ErrScanNotComplete occurs when a scan is not complete
func ErrScanNotComplete(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, ScanNotComplete, format, args...)
}

//ErrDataCannotBeEmpty occurs when the info data in a message is empty
func ErrDataCannotBeEmpty(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, DataCannotBeEmpty, format, args...)
}

//ErrDataSizeOverLimit occurs when the info data is above the size limit
func ErrDataSizeOverLimit(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, DataSizeOverLimit, format, args...)
}

//ErrCantShareWithSelf occurs when an attendee tries to share info with themselves
func ErrCantShareWithSelf(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, CantShareWithSelf, format, args...)
}

//ErrInfoAlreadyExists occurs when someone tries to share info with a person more than once
func ErrInfoAlreadyExists(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, InfoAlreadyExists, format, args...)
}

//ErrInvalidShareForScan occurs when someone tries to share info when the corresponding scan is not complete
func ErrInvalidShareForScan(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, InvalidShareForScan, format, args...)
}

//ErrDefault occurs when a random error occurs that we do not provide a unique code to
func ErrDefault(format string, args ...interface{}) sdk.Error {
	return sdk.NewError(LongyCodeSpace, DefaultError, format, args...)
}
