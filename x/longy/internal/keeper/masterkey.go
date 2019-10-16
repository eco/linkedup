package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/types"
)

// GetMasterAddress will retrieve the master key set by the module
//nolint:gocritic
func (k Keeper) GetMasterAddress(ctx sdk.Context) sdk.AccAddress {
	key := types.MasterKey()
	bz, _ := k.Get(ctx, key)

	return sdk.AccAddress(bz)
}

// SetMasterAddress will set the module's master key
//nolint:gocritic
func (k Keeper) SetMasterAddress(ctx sdk.Context, addr sdk.AccAddress) {
	key := types.MasterKey()
	bz := addr.Bytes()
	k.Set(ctx, key, bz)
}
