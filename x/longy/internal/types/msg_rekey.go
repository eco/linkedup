package types

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = MsgRekey{}

// MsgRekey implements the `sdk.Msg` interface
type MsgRekey struct {
	AttendeeID           string
	NewAttendeePublicKey sdk.AccAddress

	// expected signer of this message
	MasterPublicKey sdk.AccAddress
}

// NewMsgRekey is the creator for `RekeyMsg`
func NewMsgRekey(id string, newPublicKey, masterPublicKey sdk.AccAddress) MsgRekey {
	return MsgRekey{
		AttendeeID:           id,
		NewAttendeePublicKey: newPublicKey,
		MasterPublicKey:      masterPublicKey,
	}
}

// Route defines the route for this message
func (msg MsgRekey) Route() string {
	return RouterKey
}

// Type is the message type
func (msg MsgRekey) Type() string {
	return "rekey"
}

// ValidateBasic peforms sanity checks on the message
func (msg MsgRekey) ValidateBasic() sdk.Error {
	if msg.NewAttendeePublicKey.Empty() {
		return sdk.ErrInvalidAddress("new attendee public key is empty")
	} else if msg.MasterPublicKey.Empty() {
		return sdk.ErrInvalidAddress("master public key is empty")
	}

	return nil
}

// GetSignBytes returns the byte array that is signed over
func (msg MsgRekey) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners returns the the master public key expected to sign this message
func (msg MsgRekey) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.MasterPublicKey}
}
