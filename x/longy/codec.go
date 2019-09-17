package button

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var LONGY_ROUTE = "longy"

// RegisterCodec registers concrete types on wire codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgQrScan{}, "longy/QRScan", nil)
}
