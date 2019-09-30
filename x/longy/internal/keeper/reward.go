package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/types"
)

// GetMasterPublicKey will retrieve the reward of `kind`
func (k Keeper) GetReward(ctx sdk.Context, kind string) (types.Reward, bool) {
	key := types.RewardKey(kind)
	bz := k.Get(ctx, key)
	if bz == nil {
		return types.Reward{}, false
	}

	var r types.Reward
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &r)
	return r, true
}

// SetReward will store the `reward`
func (k Keeper) SetReward(ctx sdk.Context, reward types.Reward) {
	key := types.RewardKey(reward.Kind())
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(reward)
	k.Set(ctx, key, bz)
}
