package longy

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var LONGY_ROUTE = "longy"

// RegisterCodec registers concrete types on wire codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgQrScan{}, LONGY_ROUTE+"/QRScan", nil)
	cdc.RegisterConcrete(MsgShareInfo{}, LONGY_ROUTE+"/ShareInfo", nil)
}
