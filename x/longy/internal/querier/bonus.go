package querier

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/keeper"
)

//nolint:gocritic,unparam
func queryBonus(ctx sdk.Context, keeper keeper.Keeper) ([]byte, sdk.Error) {
	bonus := keeper.GetBonus(ctx)
	if bonus == nil {
		return nil, nil
	}

	res, err := codec.MarshalJSONIndent(keeper.Cdc, bonus)
	if err != nil {
		panic(fmt.Sprintf("json marshal bonus: %s", err))
	}

	return res, nil
}
