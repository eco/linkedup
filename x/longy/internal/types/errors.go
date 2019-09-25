package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultCodespace is the Module Name
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	QRCodeDoesNotExist sdk.CodeType = 101
)

// ErrQRCodeDoesNotExist is the error for when a QR code does not exist in our keeper
func ErrQRCodeDoesNotExist() sdk.Error {
	return sdk.NewError(DefaultCodespace, QRCodeDoesNotExist, "Name does not exist")
}
