package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

// Keeper maintains the link to data storage
type Keeper struct {
	contextStoreKey sdk.StoreKey

	accountKeeper auth.AccountKeeper
	cdc           *codec.Codec
}

// NewKeeper is a creator for `Keeper`
func NewKeeper(cdc *codec.Codec, longyStoreKey sdk.StoreKey, accKeeper auth.AccountKeeper) Keeper {
	return Keeper{
		contextStoreKey: longyStoreKey,
		accountKeeper:   accKeeper,
		cdc:             cdc,
	}
}

// AccountKeeper returns the auth module's account keeper composed with this module
func (k Keeper) AccountKeeper() auth.AccountKeeper {
	return k.accountKeeper
}

// Set sets the key value pair in the store
func (k Keeper) Set(ctx sdk.Context, key []byte, value []byte) {
	store := ctx.KVStore(k.contextStoreKey)
	store.Set(key, value)
}

// Get returns the value for the provided key from the store
func (k Keeper) Get(ctx sdk.Context, key []byte) []byte {
	store := ctx.KVStore(k.contextStoreKey)
	return store.Get(key)
}

// Delete removes the provided key value pair from the store
func (k Keeper) Delete(ctx sdk.Context, key []byte) {
	store := ctx.KVStore(k.contextStoreKey)
	store.Delete(key)
}

// Has returns whether the key exists in the store
func (k Keeper) Has(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(k.contextStoreKey)
	return store.Has(key)
}

// KVStore returns the key value store
func (k Keeper) KVStore(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.contextStoreKey)
}
