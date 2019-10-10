package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/types"
)

//GetPrize returns the prize by its id. Returns an error if it cannot find the prize with that id
//nolint:gocritic
func (k Keeper) GetPrize(ctx sdk.Context, id []byte) (prize types.GenesisPrize, err sdk.Error) {
	bz, e := k.Get(ctx, id)
	if e != nil {
		if e.Code() == types.ItemNotFound {
			err = types.ErrPrizeNotFound("invalid key passed for prize %s", id)
			return
		}
		err = e
		return
	}

	k.cdc.MustUnmarshalBinaryBare(bz, &prize)
	return
}

//GetPrizes returns all the prizes
//nolint:gocritic
func (k Keeper) GetPrizes(ctx sdk.Context) (types.GenesisPrizes, sdk.Error) {
	pz := types.GetGenesisPrizes()
	prizes := make(types.GenesisPrizes, len(pz))
	for i, p := range pz {
		prize, err := k.GetPrize(ctx, p.GetID())
		if err != nil {
			return nil, err
		}
		prizes[i] = prize
	}
	return prizes, nil
}

//SetPrize puts the prize into the store with its tier turned into the is key
//nolint:gocritic
func (k Keeper) SetPrize(ctx sdk.Context, prize *types.GenesisPrize) {
	k.Set(ctx, prize.GetID(), k.cdc.MustMarshalBinaryBare(*prize))
}
