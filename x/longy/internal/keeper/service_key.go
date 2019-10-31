package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eco/longy/x/longy/internal/types"
)

// IsServiceAccount returns true if the the account passed in is service address
//nolint:gocritic
func (k *Keeper) IsServiceAccount(ctx sdk.Context, addr sdk.Address) bool {
	key := types.ServiceKey()
	return k.isServiceAccount(ctx, addr, key)
}

// IsBonusServiceAccount returns true if the the account passed in is service address
//nolint:gocritic
func (k *Keeper) IsBonusServiceAccount(ctx sdk.Context, addr sdk.Address) bool {
	key := types.BonusServiceKey()
	return k.isServiceAccount(ctx, addr, key)
}

//IsClaimServiceAccount returns true if the the account passed in is claim service address
//nolint:gocritic
func (k *Keeper) IsClaimServiceAccount(ctx sdk.Context, addr sdk.Address) bool {
	key := types.ClaimServiceKey()
	return k.isServiceAccount(ctx, addr, key)
}

// GetService retrieves the service account and returns it
//nolint:gocritic
func (k *Keeper) GetService(ctx sdk.Context) types.GenesisService {
	key := types.ServiceKey()
	return k.getServiceAccount(ctx, key)
}

// GetBonusService retrieves the service account and returns it
//nolint:gocritic
func (k *Keeper) GetBonusService(ctx sdk.Context) types.GenesisService {
	key := types.BonusServiceKey()
	return k.getServiceAccount(ctx, key)
}

//GetClaimService retrieves the claim service account
//nolint:gocritic
func (k *Keeper) GetClaimService(ctx sdk.Context) types.GenesisService {
	key := types.ClaimServiceKey()
	return k.getServiceAccount(ctx, key)
}

// SetServiceAddress sets the service account from the genesis file
//nolint:gocritic
func (k *Keeper) SetServiceAddress(ctx sdk.Context, addr sdk.AccAddress) sdk.Error {
	key := types.ServiceKey()
	return k.setServiceAddress(ctx, addr, key)

}

// SetBonusServiceAddress -
//nolint:gocritic
func (k *Keeper) SetBonusServiceAddress(ctx sdk.Context, addr sdk.AccAddress) sdk.Error {
	key := types.BonusServiceKey()
	return k.setServiceAddress(ctx, addr, key)
}

//SetClaimServiceAddress sets the claim account for claiming prizes from the genesis file
//nolint:gocritic
func (k *Keeper) SetClaimServiceAddress(ctx sdk.Context, addr sdk.AccAddress) sdk.Error {
	key := types.ClaimServiceKey()
	return k.setServiceAddress(ctx, addr, key)
}

//nolint:gocritic
func (k *Keeper) getServiceAccount(ctx sdk.Context, key []byte) types.GenesisService {
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

//nolint:gocritic
func (k *Keeper) setServiceAddress(ctx sdk.Context, addr sdk.AccAddress, key []byte) sdk.Error {
	if addr.Empty() {
		return sdk.ErrInvalidAddress(addr.String())
	}

	account := k.accountKeeper.GetAccount(ctx, addr)
	if account == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", addr))
	}

	bz := addr.Bytes()
	k.Set(ctx, key, bz)

	return nil
}

//nolint:gocritic
func (k *Keeper) isServiceAccount(ctx sdk.Context, addr sdk.Address, key []byte) bool {
	bz, err := k.Get(ctx, key)
	if err != nil {
		return false
	}
	service := sdk.AccAddress(bz)
	return service.Equals(addr)
}
