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
	cdc.RegisterConcrete(MsgQrScan{}, RouterKey+"/QRScan", nil)
	cdc.RegisterConcrete(MsgShareInfo{}, RouterKey+"/ShareInfo", nil)
}
