package longy

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgShareInfo defines the message for sharing our info with another attendee
type MsgShareInfo struct {
	Sender    sdk.AccAddress `json:"sender"`    //Standard for all messages
	ScannedQR string         `json:"scannedQR"` //the string representation of the other attendee's QR badge
}

// NewMsgShareInfo is the constructor function for MsgShareInfo
func NewMsgShareInfo(sender sdk.AccAddress, qrCode string) MsgShareInfo {
	return MsgShareInfo{
		Sender:    sender,
		ScannedQR: qrCode,
	}
}

// Route string for this message
func (msg MsgShareInfo) Route() string { return LongyRoute }

// Type returns the message type, used to tagging transactions
func (msg MsgShareInfo) Type() string { return "qr_scan" }

// ValidateBasic performs basic checks of the message
func (msg MsgShareInfo) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}
	return nil
}

// GetSignBytes returns byte representation of the message
func (msg MsgShareInfo) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners returns the signers of this message for the authentication module
func (msg MsgShareInfo) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
