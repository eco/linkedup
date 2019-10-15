package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = MsgBonus{}
var _ sdk.Msg = MsgClearBonus{}

/** MsgBonus **/

// MsgBonus triggers a bonus period
type MsgBonus struct {
	MasterAddress sdk.AccAddress `json:"masterAddress"`
	Multiplier    uint           `json:"multiplier"`
}

// NewMsgBonus -
func NewMsgBonus(multiplier uint, masterAddr sdk.AccAddress) MsgBonus {
	return MsgBonus{
		MasterAddress: masterAddr,
		Multiplier:    multiplier,
	}
}

// Route -
func (msg MsgBonus) Route() string {
	return RouterKey
}

// Type -
func (msg MsgBonus) Type() string {
	return "bonus"
}

// ValidateBasic -
func (msg MsgBonus) ValidateBasic() sdk.Error {
	switch {
	case msg.MasterAddress.Empty():
		return sdk.ErrInvalidAddress("unset master address")
	case msg.Multiplier == 0:
		return ErrDefault("zero multiplier")
	}

	return nil
}

// GetSigners -
func (msg MsgBonus) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.MasterAddress}
}

// GetSignBytes -
func (msg MsgBonus) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

/** MsgClearBonus **/

// MsgClearBonus -
type MsgClearBonus struct {
	MasterAddress sdk.AccAddress `json:"masterAddress"`
}

// NewMsgClearBonus -
func NewMsgClearBonus(masterAddr sdk.AccAddress) MsgClearBonus {
	return MsgClearBonus{
		MasterAddress: masterAddr,
	}
}

// Route -
func (msg MsgClearBonus) Route() string {
	return RouterKey
}

// Type -
func (msg MsgClearBonus) Type() string {
	return "clear_bonus"
}

// ValidateBasic -
func (msg MsgClearBonus) ValidateBasic() sdk.Error {
	if msg.MasterAddress.Empty() {
		return sdk.ErrInvalidAddress("empty master address")
	}

	return nil
}

// GetSigners -
func (msg MsgClearBonus) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.MasterAddress}
}

// GetSignBytes -
func (msg MsgClearBonus) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
