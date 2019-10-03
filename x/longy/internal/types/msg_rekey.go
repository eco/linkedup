package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/util"
	tmcrypto "github.com/tendermint/tendermint/crypto"
)

var _ sdk.Msg = MsgRekey{}

// MsgRekey implements the `sdk.Msg` interface
type MsgRekey struct {
	AttendeeAddress      sdk.AccAddress
	NewAttendeePublicKey tmcrypto.PubKey

	// expected signer of this message
	MasterAddress sdk.AccAddress

	// Commitmentment that needs to be reclaimed
	Commitment util.Commitment
}

// NewMsgRekey is the creator for `RekeyMsg`
func NewMsgRekey(attendeeAddress, masterAddress sdk.AccAddress,
	newPublicKey tmcrypto.PubKey, commit util.Commitment) MsgRekey {

	return MsgRekey{
		AttendeeAddress:      attendeeAddress,
		NewAttendeePublicKey: newPublicKey,

		MasterAddress: masterAddress,

		Commitment: commit,
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
	if msg.AttendeeAddress.Empty() {
		return sdk.ErrInvalidAddress("attendee address is empty")
	} else if len(msg.NewAttendeePublicKey.Bytes()) == 0 {
		return sdk.ErrInvalidAddress("new attendee public key is empty")
	} else if msg.MasterAddress.Empty() {
		return sdk.ErrInvalidAddress("master public key is empty")
	}

	return nil
}

// GetSignBytes returns the byte array that is signed over
func (msg MsgRekey) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the the master public key expected to sign this message
func (msg MsgRekey) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.MasterAddress}
}
