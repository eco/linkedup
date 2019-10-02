package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc is the codec for the module
var ModuleCdc = codec.New()

func init() {
	codec.RegisterCrypto(ModuleCdc)
	RegisterCodec(ModuleCdc)
}

// RegisterCodec registers concrete types used by this module
func RegisterCodec(cdc *codec.Codec) {
	// register msgs
	cdc.RegisterConcrete(MsgQrScan{}, RouterKey+"/MsgQRScan", nil)
	cdc.RegisterConcrete(MsgShareInfo{}, RouterKey+"/MsgShareInfo", nil)
	cdc.RegisterConcrete(MsgRekey{}, RouterKey+"/MsgRekey", nil)
	cdc.RegisterConcrete(MsgClaimKey{}, RouterKey+"/MsgClaimKey", nil)

	// register types
	cdc.RegisterConcrete(Attendee{}, RouterKey+"/attendee", nil)
}
