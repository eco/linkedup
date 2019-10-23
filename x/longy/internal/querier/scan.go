package querier

import (
	"encoding/hex"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/keeper"
	"github.com/eco/longy/x/longy/internal/types"
)

//nolint:gocritic,unparam
func queryScans(ctx sdk.Context, path []string, keeper keeper.Keeper) (res []byte, err sdk.Error) {
	scan, err := keeper.GetScanByID(ctx, types.Decode(path[0]))
	if err != nil {
		return
	}
	dst := make([]byte, hex.EncodedLen(len(scan.ID)))
	hex.Encode(dst, scan.ID)
	scan.ID = dst
	res, e := codec.MarshalJSONIndent(keeper.Cdc, scan)
	if e != nil {
		panic("could not marshal result to JSON")
	}
	return
}
