package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/types"
)

// SetBonus -
func (k Keeper) SetBonus(ctx sdk.Context, b types.Bonus) {
	key := types.BonusKey()

	bz, err := k.cdc.MarshalBinaryLengthPrefixed(b)
	if err != nil {
		panic(err)
	}

	k.Set(ctx, key, bz)
}

// GetBonus -
func (k Keeper) GetBonus(ctx sdk.Context) *types.Bonus {
	key := types.BonusKey()
	bz, _ := k.Get(ctx, key)
	if bz == nil {
		return nil
	}

	var b types.Bonus
	err := k.cdc.UnmarshalBinaryLengthPrefixed(bz, &b)
	if err != nil {
		panic(err)
	}

	return &b
}

// ClearBonus -
func (k Keeper) ClearBonus(ctx sdk.Context) {
	key := types.BonusKey()
	k.Delete(ctx, key)
}

// HasBonus -
func (k Keeper) HasLiveBonus(ctx sdk.Context) bool {
	key := types.BonusKey()
	return k.Has(ctx, key)
}
