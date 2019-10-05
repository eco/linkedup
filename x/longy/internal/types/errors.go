package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultCodespace is the Module Name
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	QRCodeDoesNotExist sdk.CodeType = 101

	AttendeeDoesNotExist   sdk.CodeType = 102
	AttendeeUnclaimed      sdk.CodeType = 103
	AttendeeAlreadyClaimed sdk.CodeType = 104

	InvalidCommitmentReveal sdk.CodeType = 105
)

// ErrQRCodeDoesNotExist is the error for when a QR code does not exist in our keeper
func ErrQRCodeDoesNotExist() sdk.Error {
	return sdk.NewError(DefaultCodespace, QRCodeDoesNotExist, "name does not exist")
}

// ErrAttendeeDoesNotExist indicates an id with not corresponding attendee
func ErrAttendeeDoesNotExist() sdk.Error {
	return sdk.NewError(DefaultCodespace, AttendeeDoesNotExist, "attendee does not exist")
}

// ErrAttendeeUnclaimed indicates a attendee that is unclaimed
func ErrAttendeeUnclaimed() sdk.Error {
	return sdk.NewError(DefaultCodespace, AttendeeUnclaimed, "attendee unclaimed")
}

// ErrAttendeeAlreadyClaimed indicates a attendee that is unclaimed
func ErrAttendeeAlreadyClaimed() sdk.Error {
	return sdk.NewError(DefaultCodespace, AttendeeAlreadyClaimed, "attendee claimed")
}

// ErrInvalidCommitmentReveal indicates that the reveal is incorrect for the commitment
func ErrInvalidCommitmentReveal() sdk.Error {
	return sdk.NewError(DefaultCodespace, InvalidCommitmentReveal, "reveal to the commitment is incorrect")
}
