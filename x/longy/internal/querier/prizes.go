package querier

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/keeper"
)

//nolint:gocritic,unparam
func queryPrizes(ctx sdk.Context, keeper keeper.Keeper) (res []byte, err sdk.Error) {
	prizes, err := keeper.GetPrizes(ctx)
	if err != nil {
		return
	}

	res, e := codec.MarshalJSONIndent(keeper.Cdc, prizes)
	if e != nil {
		panic("could not marshal result to JSON")
	}

	return
}
