package longy

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgQrScan defines the message for starting off a QR scan meet of another attendee
type MsgQrScan struct {
	Sender    sdk.AccAddress `json:"sender"`    //Standard for all messages
	ScannedQR string         `json:"scannedQR"` //the string representation of the other attendee's QR badge
	ScanCode  string         `json:"scanCode"`  //the scan code from the QR account, used to validate some interaction to prevent social media posts
}

// NewMsgQrScan is the constructor function for MsgQrScan
func NewMsgQrScan(sender sdk.AccAddress, qrCode string) MsgQrScan {
	return MsgQrScan{
		Sender:    sender,
		ScannedQR: qrCode,
	}
}

// Route string for this message
func (msg MsgQrScan) Route() string { return LONGY_ROUTE }

// Type returns the message type, used to tagging transactions
func (msg MsgQrScan) Type() string { return "qr_scan" }

// ValidateBasic performs basic checks of the message
func (msg MsgQrScan) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}
	return nil
}

// GetSignBytes returns byte representation of the message
func (msg MsgQrScan) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners returns the signers of this message for the authentication module
func (msg MsgQrScan) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
