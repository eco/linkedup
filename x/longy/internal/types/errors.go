package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultCodespace is the Module Name
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	QRCodeDoesNotExist   sdk.CodeType = 101
	AttendeeDoesNotExist sdk.CodeType = 102
	RewardDoesNotExist   sdk.CodeType = 103
)

// ErrQRCodeDoesNotExist is the error for when a QR code does not exist in our keeper
func ErrQRCodeDoesNotExist() sdk.Error {
	return sdk.NewError(DefaultCodespace, QRCodeDoesNotExist, "name does not exist")
}

// ErrAttendeeDoesNotExist indicates an id with not corresponding attendee
func ErrAttendeeDoesNotExist() sdk.Error {
	return sdk.NewError(DefaultCodespace, AttendeeDoesNotExist, "attendee does not exist")
}

// ErrRewardDoesNotExist indicates a reward type that doesn't exist
func ErrRewardDoesNotExist() sdk.Error {
	return sdk.NewError(DefaultCodespace, RewardDoesNotExist, "reward does note exist")
}
