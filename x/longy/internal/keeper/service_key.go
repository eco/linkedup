package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/types"
)

//IsServiceAccount returns true if the the account passed in is service address
//nolint:gocritic
func (k *Keeper) IsServiceAccount(ctx sdk.Context, addr sdk.Address) bool {
	key := types.ServiceKey()
	bz, err := k.Get(ctx, key)
	if err != nil {
		return false
	}
	service := sdk.AccAddress(bz)
	return service.Equals(addr)
}

//GetService retrieves the service account and returns it
//nolint:gocritic
func (k *Keeper) GetService(ctx sdk.Context) types.GenesisService {
	key := types.ServiceKey()
	bz, err := k.Get(ctx, key)
	if err != nil {
		return types.GenesisService{}
	}
	address := sdk.AccAddress(bz)

	serviceAccount := k.accountKeeper.GetAccount(ctx, address)
	if serviceAccount == nil {
		return types.GenesisService{}

	}

	return types.GenesisService{
		Address: serviceAccount.GetAddress(),
		PubKey:  serviceAccount.GetPubKey(),
	}
}

//SetServiceAddress sets the service account from the genesis file
//nolint:gocritic
func (k *Keeper) SetServiceAddress(ctx sdk.Context, addr sdk.AccAddress) sdk.Error {
	if addr.Empty() {
		return sdk.ErrInvalidAddress(addr.String())
	}

	account := k.accountKeeper.GetAccount(ctx, addr)
	if account == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", addr))
	}

	key := types.ServiceKey()
	bz := addr.Bytes()
	k.Set(ctx, key, bz)

	return nil
}
