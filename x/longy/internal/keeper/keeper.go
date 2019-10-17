package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/eco/longy/x/longy/internal/types"
)

// Keeper maintains the link to data storage
type Keeper struct {
	contextStoreKey sdk.StoreKey

	accountKeeper auth.AccountKeeper
	coinKeeper    bank.Keeper
	cdc           *codec.Codec
}

// NewKeeper is a creator for `Keeper`
//nolint:gocritic
func NewKeeper(cdc *codec.Codec, longyStoreKey sdk.StoreKey, accKeeper auth.AccountKeeper,
	coinKeeper bank.Keeper) Keeper {
	return Keeper{
		contextStoreKey: longyStoreKey,
		accountKeeper:   accKeeper,
		coinKeeper:      coinKeeper,
		cdc:             cdc,
	}
}

// AccountKeeper returns the auth module's account keeper composed with this module
//nolint: gocritic
func (k Keeper) AccountKeeper() auth.AccountKeeper {
	return k.accountKeeper
}

// CoinKeeper returns the module's bank keeper
//nolint: gocritic
func (k Keeper) CoinKeeper() bank.Keeper {
	return k.coinKeeper
}

// Set sets the key value pair in the store
//nolint:gocritic
func (k Keeper) Set(ctx sdk.Context, key []byte, value []byte) {
	store := ctx.KVStore(k.contextStoreKey)
	store.Set(key, value)
}

// Get returns the value for the provided key from the store
//nolint:gocritic
func (k Keeper) Get(ctx sdk.Context, key []byte) (v []byte, err sdk.Error) {
	store := ctx.KVStore(k.contextStoreKey)
	//return store.Get(key)
	v = store.Get(key)
	if len(v) == 0 {
		err = types.ErrItemNotFound("invalid key passed for item %s", key)
		return
	}

	return
}

// Delete removes the provided key value pair from the store
//nolint:gocritic
func (k Keeper) Delete(ctx sdk.Context, key []byte) {
	store := ctx.KVStore(k.contextStoreKey)
	store.Delete(key)
}

// Has returns whether the key exists in the store
//nolint:gocritic
func (k Keeper) Has(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(k.contextStoreKey)
	return store.Has(key)
}

// KVStore returns the key value store
//nolint:gocritic
func (k Keeper) KVStore(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.contextStoreKey)
}
