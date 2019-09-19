package longy

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

//LongyRoute is the package route
var LongyRoute = "longy"

// RegisterCodec registers concrete types on wire codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgQrScan{}, LongyRoute+"/QRScan", nil)
	cdc.RegisterConcrete(MsgShareInfo{}, LongyRoute+"/ShareInfo", nil)
}
