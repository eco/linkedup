package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
)

var _ sdk.Msg = MsgBonus{}
var _ sdk.Msg = MsgClearBonus{}

/** MsgBonus **/

// MsgBonus triggers a bonus period
type MsgBonus struct {
	BonusServiceAddress sdk.AccAddress `json:"bonus_service_address"`
	Multiplier          string         `json:"multiplier"`
}

// NewMsgBonus -
func NewMsgBonus(multiplier string, addr sdk.AccAddress) MsgBonus {
	return MsgBonus{
		BonusServiceAddress: addr,
		Multiplier:          multiplier,
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
	case msg.BonusServiceAddress.Empty():
		return sdk.ErrInvalidAddress("unset bonus service address")
	case len(msg.Multiplier) == 0:
		return ErrDefault("zero multiplier")
	}

	num, err := strconv.ParseFloat(msg.Multiplier, 64)
	if err != nil {
		return ErrDefault("non-parsable multiplier")
	} else if num <= 0 {
		return ErrDefault("multiplier must be >= 0")
	}

	return nil
}

// GetSigners -
func (msg MsgBonus) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.BonusServiceAddress}
}

// GetSignBytes -
func (msg MsgBonus) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

/** MsgClearBonus **/

// MsgClearBonus -
type MsgClearBonus struct {
	BonusServiceAddress sdk.AccAddress `json:"bonus_service_address"`
}

// NewMsgClearBonus -
func NewMsgClearBonus(addr sdk.AccAddress) MsgClearBonus {
	return MsgClearBonus{
		BonusServiceAddress: addr,
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
	if msg.BonusServiceAddress.Empty() {
		return sdk.ErrInvalidAddress("empty bonus service address")
	}

	return nil
}

// GetSigners -
func (msg MsgClearBonus) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.BonusServiceAddress}
}

// GetSignBytes -
func (msg MsgClearBonus) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
