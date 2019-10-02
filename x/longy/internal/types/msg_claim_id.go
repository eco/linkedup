package types

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgClaimID defines the message for assigning an address to an id
type MsgClaimID struct {
	Sender  sdk.AccAddress `json:"sender"`  //Standard for all messages
	ID      string         `json:"id"`      //the string representation of the attendee's QR badgeID
	Address sdk.AccAddress `json:"address"` //the account address to be associated to the id
}

// NewMsgClaimID is the constructor function for MsgClaimID
func NewMsgClaimID(sender sdk.AccAddress, id string, address sdk.AccAddress) MsgClaimID {
	return MsgClaimID{
		Sender:  sender,
		ID:      id,
		Address: address,
	}
}

// Route string for this message
func (msg MsgClaimID) Route() string { return RouterKey }

// Type returns the message type, used to tagging transactions
func (msg MsgClaimID) Type() string { return "claim_id" }

// ValidateBasic performs basic checks of the message
func (msg MsgClaimID) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}
	return nil
}

// GetSignBytes returns byte representation of the message
func (msg MsgClaimID) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners returns the signers of this message for the authentication module
func (msg MsgClaimID) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
