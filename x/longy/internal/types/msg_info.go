package types

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	//MaxDataSize max payload for data
	MaxDataSize = 2048 //2k
)

// MsgInfo defines the message for sharing our info with another attendee
type MsgInfo struct {
	Sender   sdk.AccAddress `json:"sender"`   //Standard for all messages
	Receiver sdk.AccAddress `json:"receiver"` //the person that is getting the info
	Data     []byte         `json:"data"`     //the encrypted data to store
}

// NewMsgInfo is the constructor function for MsgInfo
func NewMsgInfo(sender sdk.AccAddress, receiver sdk.AccAddress, data []byte) MsgInfo {
	return MsgInfo{
		Sender:   sender,
		Receiver: receiver,
		Data:     data,
	}
}

// Route string for this message
func (msg MsgInfo) Route() string { return RouterKey }

// Type returns the message type, used to tagging transactions
func (msg MsgInfo) Type() string { return "info" }

// ValidateBasic performs basic checks of the message
func (msg MsgInfo) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}

	if msg.Receiver.Empty() {
		return sdk.ErrInvalidAddress(msg.Receiver.String())
	}

	if msg.Receiver.Equals(msg.Sender) {
		return ErrCantShareWithSelf("cannot share info with self")
	}

	if len(msg.Data) == 0 {
		return ErrDataCannotBeEmpty("info data cannot be empty")
	}

	if len(msg.Data) > MaxDataSize {
		return ErrDataSizeOverLimit("info data size is over the limit of %d bytes", MaxDataSize)
	}
	return nil
}

// GetSignBytes returns byte representation of the message
func (msg MsgInfo) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners returns the signers of this message for the authentication module
func (msg MsgInfo) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
