package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
	tmcrypto "github.com/tendermint/tendermint/crypto"
)

var _ sdk.Msg = MsgKey{}

// MsgKey implements the `sdk.Msg` interface
type MsgKey struct {
	AttendeeAddress      sdk.AccAddress  `json:"attendeeAddress"`
	NewAttendeePublicKey tmcrypto.PubKey `json:"newAttendeePublicKey"`

	// expected signer of this message
	MasterAddress sdk.AccAddress `json:"masterAddress"`

	// Commitmentment that needs to be reclaimed
	Commitment util.Commitment `json:"commitment"`
}

// NewMsgKey is the creator for `RekeyMsg`
func NewMsgKey(attendeeAddress, masterAddress sdk.AccAddress,
	newPublicKey tmcrypto.PubKey, commit util.Commitment) MsgKey {

	return MsgKey{
		AttendeeAddress:      attendeeAddress,
		NewAttendeePublicKey: newPublicKey,

		MasterAddress: masterAddress,

		Commitment: commit,
	}
}

// Route defines the route for this message
//nolint:gocritic
func (msg MsgKey) Route() string {
	return RouterKey
}

// Type is the message type
//nolint:gocritic
func (msg MsgKey) Type() string {
	return "key"
}

// ValidateBasic peforms sanity checks on the message
//nolint:gocritic
func (msg MsgKey) ValidateBasic() sdk.Error {
	if msg.AttendeeAddress.Empty() {
		return sdk.ErrInvalidAddress("attendee address is empty")
	} else if len(msg.NewAttendeePublicKey.Bytes()) == 0 {
		return sdk.ErrInvalidAddress("new attendee public key is empty")
	}

	if msg.MasterAddress.Empty() {
		return sdk.ErrInvalidAddress("master public key is empty")
	}

	return nil
}

// GetSignBytes returns the byte array that is signed over
//nolint:gocritic
func (msg MsgKey) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the the master public key expected to sign this message
//nolint:gocritic
func (msg MsgKey) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.MasterAddress}
}
