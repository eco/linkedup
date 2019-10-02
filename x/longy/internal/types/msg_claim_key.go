package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = MsgClaimKey{}

// MsgClaimKey is used to claim a prior rekey message
type MsgClaimKey struct {
	AttendeeAddress sdk.AccAddress
	Secret          []byte

	Signer sdk.AccAddress
}

// NewMsgClaimKey in the constructor for `MsgClaimKey`
func NewMsgClaimKey(address, signer sdk.AccAddress, secret []byte) MsgClaimKey {
	return MsgClaimKey{
		AttendeeAddress: address,
		Secret:          secret,
		Signer:          signer,
	}
}

// Route defines the route for this message
func (msg MsgClaimKey) Route() string {
	return RouterKey
}

// Type is the message type
func (msg MsgClaimKey) Type() string {
	return "claimkey"
}

// ValidateBasic performs sanity checks on the message
func (msg MsgClaimKey) ValidateBasic() sdk.Error {
	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress("empty signer")
	} else if msg.AttendeeAddress.Empty() {
		return sdk.ErrInvalidAddress("empty attendee address")
	} else if len(msg.Secret) == 0 {
		return sdk.ErrNoSignatures("missing secret")
	}

	return nil
}

// GetSignBytes returns the byte array that is signed over
func (msg MsgClaimKey) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the
func (msg MsgClaimKey) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
