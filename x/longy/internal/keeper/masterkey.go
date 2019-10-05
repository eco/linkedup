package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/types"
)

// GetMasterPublicKey will retrieve the master key set by the module
//nolint:gocritic
func (k Keeper) GetMasterPublicKey(ctx sdk.Context) sdk.AccAddress {
	key := types.MasterKey()
	bz, _ := k.Get(ctx, key)

	return sdk.AccAddress(bz)
}

// SetMasterPublicKey will set the module's master key
//nolint:gocritic
func (k Keeper) SetMasterPublicKey(ctx sdk.Context, publicKey sdk.AccAddress) {
	key := types.MasterKey()
	bz := publicKey.Bytes()
	k.Set(ctx, key, bz)
}
