package types

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = RekeyMsg{}

// RekeyMsg implements the `sdk.Msg` interface
type RekeyMsg struct {
	AttendeeID         int
	MasterKeySignature []byte
}

// NewRekeyMsg is the creator for `RekeyMsg`
func NewRekeyMsg(id int, masterSig []byte) RekeyMsg {
	return RekeyMsg{
		AttendeeID:         id,
		MasterKeySignature: masterSig,
	}
}

// Route defines the route for this message
func (m RekeyMsg) Route() string {
	return RouterKey
}

// Type is the message type
func (m RekeyMsg) Type() string {
	return "rekey"
}

// ValidateBasic peforms sanity checks on the message
func (m RekeyMsg) ValidateBasic() sdk.Error {
	if len(m.MasterKeySignature) == 0 {
		return sdk.ErrNoSignatures("master signature missing")
	}

	return nil
}

// GetSignBytes returns the byte array that is signed over
func (m RekeyMsg) GetSignBytes() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners for this message is `nil`. The only signer for this message type
// is the master public key of the game. Might change this to return it's hardcoded value
func (m RekeyMsg) GetSigners() []sdk.AccAddress {
	return nil
}
