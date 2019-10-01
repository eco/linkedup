package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// ModuleCdc is the codec for the module
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}

// RegisterCodec registers concrete types on wire codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgClaimID{}, RouterKey+"/claimId", nil)
	cdc.RegisterConcrete(MsgQrScan{}, RouterKey+"/qrScan", nil)
	cdc.RegisterConcrete(MsgShareInfo{}, RouterKey+"/shareInfo", nil)
}
