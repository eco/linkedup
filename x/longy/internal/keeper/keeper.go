package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage
type Keeper struct {
	contextStoreKey sdk.StoreKey
	cdc             *codec.Codec
}

// NewKeeper is a creator for `Keeper`
func NewKeeper(scanKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		contextStoreKey: scanKey,
		cdc:             cdc,
	}
}

// Set sets the key value pair in the store
func (ds Keeper) Set(ctx sdk.Context, key []byte, value []byte) {
	store := ctx.KVStore(ds.contextStoreKey)
	store.Set(key, value)
}

// Get returns the value for the provided key from the store
func (ds Keeper) Get(ctx sdk.Context, key []byte) []byte {
	store := ctx.KVStore(ds.contextStoreKey)
	return store.Get(key)
}

// Delete removes the provided key value pair from the store
func (ds Keeper) Delete(ctx sdk.Context, key []byte) {
	store := ctx.KVStore(ds.contextStoreKey)
	store.Delete(key)
}

// Has returns whether the key exists in the store
func (ds Keeper) Has(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(ds.contextStoreKey)
	return store.Has(key)
}

// KVStore returns the key value store
func (ds Keeper) KVStore(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(ds.contextStoreKey)
}
