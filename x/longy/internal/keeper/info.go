package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/types"
)

//GetInfoByID is the getter or info by id
//nolint:gocritic
func (k Keeper) GetInfoByID(ctx sdk.Context, id []byte) (info *types.Info, err sdk.Error) {
	bz, e := k.Get(ctx, id)
	if e != nil {
		if e.Code() == types.ItemNotFound {
			err = types.ErrInfoNotFound("invalid key passed for info %s", id)
			return
		}
		err = e
		return
	}

	k.cdc.MustUnmarshalBinaryBare(bz, &info)
	return
}

//SetInfo sets the info by id
//nolint:gocritic
func (k Keeper) SetInfo(ctx sdk.Context, info *types.Info) {
	k.Set(ctx, info.ID, k.cdc.MustMarshalBinaryBare(&info))
}
