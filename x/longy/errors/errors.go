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
