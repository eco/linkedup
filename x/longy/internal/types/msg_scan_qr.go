package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
)

// MsgScanQr defines the message for starting off a QR scan meet of another attendee
type MsgScanQr struct {
	Sender    sdk.AccAddress `json:"sender"`         //Standard for all messages
	ScannedQR string         `json:"scannedQR"`      //the string representation of the other attendee's QR badge
	Data      []byte         `json:"data,omitempty"` //the encrypted data to store
	//// some interaction to prevent social media posts
	//ScanCode  string         `json:"scanCode"`  //the scan code from the QR account, used to validate
}

// NewMsgQrScan is the constructor function for MsgScanQr
func NewMsgQrScan(sender sdk.AccAddress, qrCode string, data []byte) MsgScanQr {
	return MsgScanQr{
		Sender:    sender,
		ScannedQR: qrCode,
		Data:      data,
	}
}

// Route string for this message
func (msg MsgScanQr) Route() string { return RouterKey }

// Type returns the message type, used to tagging transactions
func (msg MsgScanQr) Type() string { return "qr_scan" }

// ValidateBasic performs basic checks of the message
func (msg MsgScanQr) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}

	if !ValidQrCode(msg.ScannedQR) {
		return ErrQRCodeInvalid("message QR code is invalid, should be a string of a positive integer")
	}

	//data can me empty
	return nil
}

// GetSignBytes returns byte representation of the message
func (msg MsgScanQr) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the signers of this message for the authentication module
func (msg MsgScanQr) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

//ValidQrCode returns true if the qr code is valid, ie its a positive integer
func ValidQrCode(code string) bool {
	//valid qr code is a 10 digit integer
	val, err := strconv.Atoi(code)
	if err != nil || val <= 0 {
		return false
	}
	return true
}
