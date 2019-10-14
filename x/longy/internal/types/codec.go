package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc is the codec for the module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	codec.RegisterCrypto(ModuleCdc)
	RegisterCodec(ModuleCdc)
}

// RegisterCodec registers concrete types used by this module
func RegisterCodec(cdc *codec.Codec) {
	// register msgs
	cdc.RegisterConcrete(MsgScanQr{}, RouterKey+"/MsgScanQr", nil)
	cdc.RegisterConcrete(MsgInfo{}, RouterKey+"/MsgInfo", nil)
	cdc.RegisterConcrete(MsgKey{}, RouterKey+"/MsgKey", nil)
	cdc.RegisterConcrete(MsgClaimKey{}, RouterKey+"/MsgClaimKey", nil)
	cdc.RegisterConcrete(MsgBonus{}, RouterKey+"/MsgBonus", nil)
	cdc.RegisterConcrete(MsgClearBonus{}, RouterKey+"/MsgClearBonus", nil)

	// register types
	cdc.RegisterConcrete(Attendee{}, RouterKey+"/Attendee", nil)
	cdc.RegisterConcrete(Bonus{}, RouterKey+"/Bonus", nil)
}
