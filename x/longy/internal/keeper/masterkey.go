package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/types"
)

//IsMasterAccount returns true if the the account passed in is master
//nolint:gocritic
func (k Keeper) IsMasterAccount(ctx sdk.Context, addr sdk.Address) bool {
	key := types.MasterKey()
	bz, err := k.Get(ctx, key)
	if err != nil {
		return false
	}
	master := sdk.AccAddress(bz)
	return master.Equals(addr)
}

//SetMasterAddress sets the master account from the genesis file
//nolint:gocritic
func (k Keeper) SetMasterAddress(ctx sdk.Context, addr sdk.AccAddress) sdk.Error {
	if addr.Empty() {
		return sdk.ErrInvalidAddress(addr.String())
	}

	account := k.accountKeeper.GetAccount(ctx, addr)
	if account == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", addr))
	}

	key := types.MasterKey()
	bz := addr.Bytes()
	k.Set(ctx, key, bz)

	return nil
}
